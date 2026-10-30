package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/action"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Controller) runActions(ctx context.Context) error {
	filePaths, err := action.Find()
	if err != nil {
		return fmt.Errorf("find action files: %w", err)
	}
	failed := false
	for _, filePath := range filePaths {
		logger := c.logger.With("action_file_path", filePath)
		vw := &validateAction{
			action:  filePath,
			logger:  logger,
			gh:      c.gh,
			rootDir: c.rootDir,
		}
		if err := vw.validate(ctx); err != nil {
			slogerr.WithError(logger, err).Error("validate action")
			failed = true
		}
	}
	if failed {
		return errors.New("some action files are invalid")
	}
	return nil
}

type validateAction struct {
	action  string
	logger  *slog.Logger
	gh      GitHub
	rootDir string
}

func (v *validateAction) validate(ctx context.Context) error {
	act := &workflow.Action{}
	if err := workflow.ReadAction(v.action, act); err != nil {
		return fmt.Errorf("read an action file: %w", err)
	}
	failed := false
	for _, step := range act.Runs.Steps {
		vs := &validateStep{
			step:    step,
			logger:  v.logger,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vs.validate(ctx); err != nil {
			slogerr.WithError(v.logger, err).Error("validate a step")
			failed = true
		}
	}
	if failed {
		return errors.New("some steps are invalid")
	}
	return nil
}
