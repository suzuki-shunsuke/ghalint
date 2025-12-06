package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

func (c *Controller) runWorkflow(ctx context.Context) error {
	filePaths, err := workflow.List(c.fs)
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	failed := false
	for _, filePath := range filePaths {
		logger := c.logger.With("workflow_file_path", filePath)
		vw := &validateWorkflow{
			workflow: filePath,
			logger:   logger,
			fs:       c.fs,
			gh:       c.gh,
			rootDir:  c.rootDir,
		}
		if err := vw.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, urfave.ErrSilent) {
				slogerr.WithError(logger, err).Error("validate workflow")
			}
		}
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}

type validateWorkflow struct {
	workflow string
	logger   *slog.Logger
	fs       afero.Fs
	gh       GitHub
	rootDir  string
}

func (v *validateWorkflow) validate(ctx context.Context) error {
	wf := &workflow.Workflow{
		FilePath: v.workflow,
	}
	if err := workflow.Read(v.fs, v.workflow, wf); err != nil {
		return fmt.Errorf("read a workflow file: %w", err)
	}
	failed := false
	for name, job := range wf.Jobs {
		vj := &validateJob{
			job:     job,
			logger:  v.logger.With("job_key", name),
			fs:      v.fs,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vj.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, urfave.ErrSilent) {
				slogerr.WithError(v.logger, err).Error("validate job")
			}
		}
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}
