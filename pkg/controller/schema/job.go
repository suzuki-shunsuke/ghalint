package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type validateJob struct {
	job     *workflow.Job
	logE    *logrus.Entry
	fs      afero.Fs
	gh      GitHub
	rootDir string
}

func (v *validateJob) validate(ctx context.Context) error {
	// Get actions
	if v.job.Uses != "" {
		v.logE = v.logE.WithField("reusable_workflow", v.job.Uses)
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
			logE:    v.logE,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vs.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, SilentError) {
				logerr.WithError(v.logE, err).Error("validate a step")
			}
		}
	}
	if failed {
		return SilentError
	}
	return nil
}
