package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) Run(ctx context.Context) error {
	// Find action.yaml and workflow files
	failed := false
	if err := c.runWorkflow(ctx); err != nil {
		failed = true
		if !errors.Is(err, SilentError) {
			logerr.WithError(c.logE, err).Error("validate workflows")
		}
	}
	if err := c.runActions(ctx); err != nil {
		if !errors.Is(err, SilentError) {
			return fmt.Errorf("validate actions: %w", err)
		}
		return SilentError
	}
	if failed {
		return SilentError
	}
	return nil
}
