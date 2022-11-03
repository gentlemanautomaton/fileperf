package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
)

func main() {
	// Capture shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var cli struct {
		Scan ScanCmd `kong:"cmd,help='Scans a set of file paths recursively and tests the read performance of each file. Reports cumulative statistics.'"`
	}

	app := kong.Parse(&cli,
		kong.Description("Tests the read performance of the file system."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	if err := app.Run(ctx); err != nil {
		app.FatalIfErrorf(err)
	}
}
