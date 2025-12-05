package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/go-stdutil"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "ghalint",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&stdutil.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	}, afero.NewOsFs(), logger, logLevelVar)
	if err := app.Run(ctx, os.Args); err != nil {
		hasLogLevel := &controller.HasLogLevelError{}
		if errors.As(err, &hasLogLevel) {
			slogerr.WithError(logger, hasLogLevel.Err).Log(ctx, hasLogLevel.LogLevel, "ghalint failed")
			return 1
		}
		if errors.Is(err, schema.ErrSilent) {
			return 1
		}
		slogerr.WithError(logger, err).Error("ghalint failed")
		return 1
	}
	return 0
}
