package cli

import (
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/urfave/cli/v2"
)

func (r *Runner) Run(ctx *cli.Context) error {
	logE := r.logE

	if color := ctx.String("log-color"); color != "" {
		log.SetColor(color, logE)
	}

	if logLevel := ctx.String("log-level"); logLevel != "" {
		log.SetLevel(logLevel, logE)
	}

	ctrl := controller.New(r.fs)

	return ctrl.Run(ctx.Context, logE, ctx.String("config")) //nolint:wrapcheck
}
