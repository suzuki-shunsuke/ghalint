package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/log"
	"github.com/urfave/cli/v3"
)

func (r *Runner) Run(ctx context.Context, cmd *cli.Command) error {
	logE := r.logE

	if err := log.Set(logE, cmd.String("log-level"), cmd.String("log-color")); err != nil {
		return fmt.Errorf("configure logger: %w", err)
	}

	ctrl := controller.New(r.fs)

	return ctrl.Run(ctx, logE, cmd.String("config")) //nolint:wrapcheck
}
