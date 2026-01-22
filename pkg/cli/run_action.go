package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller/act"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

func (r *Runner) RunAction(ctx context.Context, logger *slogutil.Logger, args *RunActionArgs) error {
	if err := logger.SetColor(args.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	ctrl := act.New(r.fs)

	return ctrl.Run(ctx, logger.Logger, args.Config, args.Files...) //nolint:wrapcheck
}
