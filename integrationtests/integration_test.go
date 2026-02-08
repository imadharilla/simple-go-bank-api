package integrationtests

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"tiny-bank-api/api"
	"tiny-bank-api/pkg/database"
	"tiny-bank-api/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestMain(m *testing.M) {
	// Set up database connection for integration tests
	postgresURL := getEnvOrDefault("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/sumup_bank")

	ctx := context.Background()
	db, err := database.NewConnection(ctx, postgresURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = db.Close()
	}()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sqldb := database.LoggingDB{SQLDB: db, Logger: logger}
	s := store.NewStore(sqldb)

	testHandler = newTestService(logger, s)

	code := m.Run()
	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func newTestService(logger *slog.Logger, store store.Store) *chi.Mux {
	apiHandler := api.NewAPI(logger, store)
	apiStrictHandler := api.NewStrictHandlerWithOptions(
		apiHandler,
		nil,
		api.StrictHTTPServerOptions{},
	)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Mount("/", api.HandlerFromMux(apiStrictHandler, nil))
	})
	return router
}
