package cli

import (
	"context"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/urfave/cli/v3"
)

func (r *Runner) RunAction(ctx context.Context, cmd *cli.Command) error {
	logE := r.logE

	if color := cmd.String("log-color"); color != "" {
		log.SetColor(color, logE)
	}

	if logLevel := cmd.String("log-level"); logLevel != "" {
		log.SetLevel(logLevel, logE)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx, logE, cmd.String("config"), cmd.Args().Slice()...) //nolint:wrapcheck
}
