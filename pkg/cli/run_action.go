package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

func (r *Runner) RunAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.String("log-color") != "" {
		r.logger.Warn("log color option is deprecated and doesn't work anymore. This is kept for backward compatibility.")
	}
	if err := slogutil.SetLevel(r.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx, r.logger, cmd.String("config"), cmd.Args().Slice()...) //nolint:wrapcheck
}
