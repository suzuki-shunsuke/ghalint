package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	logE := logrus.NewEntry(logrus.New())
	if err := core(); err != nil {
		hasLogLevel := &controller.HasLogLevelError{}
		if errors.As(err, &hasLogLevel) {
			logerr.WithError(logE, hasLogLevel.Err).Log(hasLogLevel.LogLevel, "ghalint failed")
			os.Exit(1)
		}
		logerr.WithError(logE, err).Fatal("ghalint failed")
	}
}

func core() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&cli.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	}, afero.NewOsFs())
	return app.RunContext(ctx, os.Args) //nolint:wrapcheck
}
