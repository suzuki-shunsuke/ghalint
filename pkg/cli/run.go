package cli

import (
	"context"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

func (r *Runner) Run(ctx context.Context, logger *slogutil.Logger, args *RunArgs) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(args.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}

	ctrl := controller.New(r.fs)

	return ctrl.Run(ctx, logger.Logger, args.Config) //nolint:wrapcheck
}
