package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type validateJob struct {
	job     *workflow.Job
	logger  *slog.Logger
	fs      afero.Fs
	gh      GitHub
	rootDir string
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
			step:    step,
			fs:      v.fs,
			logger:  v.logger,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vs.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, ErrSilent) {
				slogerr.WithError(v.logger, err).Error("validate a step")
			}
		}
	}
	if failed {
		return ErrSilent
	}
	return nil
}
