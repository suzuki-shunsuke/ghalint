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

type validateJob struct {
	job     *workflow.Job
	logger  *slog.Logger
	gh      GitHub
	rootDir string
	// workflow is the path of the workflow file the job is in.
	workflow string
	lockfile *lockfile.Lockfile
}

func (v *validateJob) validate(ctx context.Context) error {
	// Get actions
	if v.job.Uses != "" {
		v.logger = v.logger.With("reusable_workflow", v.job.Uses)
		if err := v.validateReusableWorkflow(ctx); err != nil {
			return fmt.Errorf("validate a reusable workflow: %w", err)
		}
		return nil
	}
	failed := false
	for _, step := range v.job.Steps {
		vs := &validateStep{
			step:     step,
			logger:   v.logger,
			gh:       v.gh,
			rootDir:  v.rootDir,
			workflow: v.workflow,
			lockfile: v.lockfile,
		}
		if err := vs.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, urfave.ErrSilent) {
				slogerr.WithError(v.logger, err).Error("validate a step")
			}
		}
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}
