package policy

import (
	"errors"

	"github.com/sirupsen/logrus"
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

func (p *JobTimeoutMinutesIsRequiredPolicy) ApplyJob(_ *logrus.Entry, _ *config.Config, _ *JobContext, job *workflow.Job) error {
	if job.TimeoutMinutes != 0 {
		return nil
	}
	if job.Uses != "" {
		// when a reusable workflow is called with "uses", "timeout-minutes" is not available.
		return nil
	}
	for _, step := range job.Steps {
		if step.TimeoutMinutes == 0 {
			return errors.New("job's timeout-minutes is required")
		}
	}
	return nil
}
