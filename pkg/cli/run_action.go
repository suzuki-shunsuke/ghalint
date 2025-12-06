package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

func (r *Runner) RunAction(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	if err := logger.SetColor(cmd.String("log-color")); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx, logger.Logger, cmd.String("config"), cmd.Args().Slice()...) //nolint:wrapcheck
}
