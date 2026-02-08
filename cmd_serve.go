package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tiny-bank-api/api"
	"tiny-bank-api/pkg/database"
	"tiny-bank-api/pkg/logging"
	"tiny-bank-api/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CmdServe struct {
	ListenAddress    string `help:"Port to listen on." default:"localhost:8080" env:"LISTEN_PORT"`
	PostgresUser     string `name:"postgresuser" help:"Username to authenticate with." default:"postgresuser" env:"POSTGRES_USER"`
	PostgresPassword string `name:"postgrespassword" help:"Password to authenticate with." default:"postgresuser" env:"POSTGRES_PASSWORD"`
	PostgresHost     string `name:"postgreshost" help:"Host of the postgresql database." default:"localhost:5432" env:"POSTGRES_HOST"`
}

func (c CmdServe) Run() error {
	postgresURL := "postgres://" + c.PostgresUser + ":" + c.PostgresPassword + "@" + c.PostgresHost + "/sumup_bank"
	logger := logging.ProdLogger()

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelFunc()

	db, err := database.NewConnection(ctx, postgresURL)
	if err != nil {
		logger.Error("Error creating database connection: " + err.Error())
		return fmt.Errorf("error creating database connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing db conn: " + err.Error())
		}
	}()
	sqldb := database.LoggingDB{SQLDB: db, Logger: logger}
	s := store.NewStore(sqldb)

	svc := NewService(logger, s)

	server := &http.Server{
		Addr:    c.ListenAddress,
		Handler: svc,
	}

	go func() {
		logger.Info("Starting HTTP server", "address", c.ListenAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error starting http server: " + err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down HTTP server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("error shutting down http server: " + err.Error())
	}

	return nil
}

func NewService(logger *slog.Logger, store store.Store) *chi.Mux {

	apiHandler := api.NewAPI(logger, store)
	apiStrictHandler := api.NewStrictHandlerWithOptions(
		apiHandler,
		nil,
		api.StrictHTTPServerOptions{ResponseErrorHandlerFunc: HandleResponseError},
	)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(injectRequestIntoContext)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// No-op
	})

	router.Route("/api", func(r chi.Router) {
		r.Mount("/", api.HandlerFromMux(apiStrictHandler, nil))
	})
	return router
}

func HandleResponseError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("Failed to handle request.", "error", err, "path", r.URL.Path, "method", r.Method)

	w.WriteHeader(http.StatusInternalServerError)
}

type requestKey struct{}

func injectRequestIntoContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
