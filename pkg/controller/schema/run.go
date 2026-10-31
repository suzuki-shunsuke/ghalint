package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

func (c *Controller) Run(ctx context.Context) error {
	if err := c.readLockfile(); err != nil {
		return err
	}
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

// readLockfile reads the lockfile so that commit SHAs of locked actions are read from it instead of the GitHub API.
// If the lockfile can't be read, it only logs a warning.
func (c *Controller) readLockfile() error {
	filePaths, err := workflow.List()
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	lf, err := lockfile.Read(filePaths)
	if err != nil {
		slogerr.WithError(c.logger, err).Warn("read the lockfile. Commit SHAs of actions are resolved through the GitHub API", "lockfile", lockfile.Path)
		return nil
	}
	c.shas.lockfile = lf
	return nil
}
