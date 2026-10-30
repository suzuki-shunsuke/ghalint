package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

func (c *Controller) runWorkflow(ctx context.Context) error {
	filePaths, err := workflow.List()
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	lf, err := lockfile.Read(filePaths)
	if err != nil {
		slogerr.WithError(c.logger, err).Warn("read the lockfile. Commit SHAs of actions are resolved through the GitHub API", "lockfile", lockfile.Path)
	}
	failed := false
	for _, filePath := range filePaths {
		logger := c.logger.With("workflow_file_path", filePath)
		vw := &validateWorkflow{
			workflow: filePath,
			logger:   logger,
			gh:       c.gh,
			rootDir:  c.rootDir,
			lockfile: lf,
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
	gh       GitHub
	rootDir  string
	lockfile *lockfile.Lockfile
}

func (v *validateWorkflow) validate(ctx context.Context) error {
	wf := &workflow.Workflow{
		FilePath: v.workflow,
	}
	if err := workflow.Read(v.workflow, wf); err != nil {
		return fmt.Errorf("read a workflow file: %w", err)
	}
	failed := false
	for name, job := range wf.Jobs {
		vj := &validateJob{
			job:      job,
			logger:   v.logger.With("job_key", name),
			gh:       v.gh,
			rootDir:  v.rootDir,
			workflow: v.workflow,
			lockfile: v.lockfile,
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
