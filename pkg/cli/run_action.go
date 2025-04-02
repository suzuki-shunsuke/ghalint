package cli

import (
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/urfave/cli/v3"
)

func (r *Runner) RunAction(ctx context.Context, ctx *cli.Command) error {
	logE := r.logE

	if color := ctx.String("log-color"); color != "" {
		log.SetColor(color, logE)
	}

	if logLevel := ctx.String("log-level"); logLevel != "" {
		log.SetLevel(logLevel, logE)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx.Context, logE, ctx.String("config"), ctx.Args().Slice()...) //nolint:wrapcheck
}
