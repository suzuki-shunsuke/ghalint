package policy

import (
	"errors"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type JobTimeoutMinutesIsRequiredPolicy struct{}

func (p *JobTimeoutMinutesIsRequiredPolicy) Name() string {
	return "job_timeout_minutes_is_required"
}

func (p *JobTimeoutMinutesIsRequiredPolicy) ID() string {
	return "012"
}

func (p *JobTimeoutMinutesIsRequiredPolicy) ApplyJob(_ *slog.Logger, _ *config.Config, _ *JobContext, job *workflow.Job) error {
	if job.TimeoutMinutes != nil {
		return nil
	}
	if job.Uses != "" {
		// when a reusable workflow is called with "uses", "timeout-minutes" is not available.
		return nil
	}
	for _, step := range job.Steps {
		if step.TimeoutMinutes == nil {
			return errors.New("job's timeout-minutes is required")
		}
	}
	return nil
}
