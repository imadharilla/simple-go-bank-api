package main

import (
	"log/slog"
	"tiny-bank-api/pkg/logging"

	"github.com/alecthomas/kong"
)

type Cli struct {
	Serve CmdServe `cmd:"1" help:"Run the API to serve requests."`
}

func main() {
	var cli Cli
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Error)

	logger := logging.ProdLogger()
	slog.SetDefault(logging.ProdLogger())

	err := ctx.Run()
	if err != nil {
		logger.Error("error in main: " + err.Error())
	}
}
