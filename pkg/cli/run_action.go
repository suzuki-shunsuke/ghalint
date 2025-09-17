package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/logrus-util/log"
	"github.com/urfave/cli/v3"
)

func (r *Runner) RunAction(ctx context.Context, cmd *cli.Command) error {
	logE := r.logE

	if err := log.Set(logE, cmd.String("log-level"), cmd.String("log-color")); err != nil {
		return fmt.Errorf("configure logger: %w", err)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx, logE, cmd.String("config"), cmd.Args().Slice()...) //nolint:wrapcheck
}
