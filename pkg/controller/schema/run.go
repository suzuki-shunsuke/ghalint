package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

func (c *Controller) Run(ctx context.Context) error {
	// Find action.yaml and workflow files
	failed := false
	if err := c.runWorkflow(ctx); err != nil {
		failed = true
		if !errors.Is(err, urfave.ErrSilent) {
			slogerr.WithError(c.logger, err).Error("validate workflows")
		}
	}
	if err := c.runActions(ctx); err != nil {
		if !errors.Is(err, urfave.ErrSilent) {
			return fmt.Errorf("validate actions: %w", err)
		}
		return urfave.ErrSilent
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}
