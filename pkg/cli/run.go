package cli

import (
	"os"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/urfave/cli/v2"
)

func (r *Runner) Run(ctx *cli.Context) error {
	logE := log.New(r.flags.Version)

	if color := os.Getenv("GHALINT_LOG_COLOR"); color != "" {
		log.SetColor(color, logE)
	}

	ctrl := controller.New(r.fs)

	return ctrl.Run(ctx.Context, logE) //nolint:wrapcheck
}
