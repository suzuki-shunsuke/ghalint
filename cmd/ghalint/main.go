package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&cli.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	return app.RunContext(ctx, os.Args) //nolint:wrapcheck
}
