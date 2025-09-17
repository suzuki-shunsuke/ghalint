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
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/go-stdutil"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/suzuki-shunsuke/logrus-util/log"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	logE := log.New("ghalint", version)
	if err := core(logE); err != nil {
		hasLogLevel := &controller.HasLogLevelError{}
		if errors.As(err, &hasLogLevel) {
			logerr.WithError(logE, hasLogLevel.Err).Log(hasLogLevel.LogLevel, "ghalint failed")
			os.Exit(1)
		}
		if errors.Is(err, schema.ErrSilent) {
			os.Exit(1)
		}
		logerr.WithError(logE, err).Fatal("ghalint failed")
	}
}

func core(logE *logrus.Entry) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&stdutil.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	}, afero.NewOsFs(), logE)
	return app.Run(ctx, os.Args) //nolint:wrapcheck
}
