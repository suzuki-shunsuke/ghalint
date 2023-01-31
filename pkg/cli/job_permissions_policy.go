package cli

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type JobPermissionsPolicy struct{}

func (policy *JobPermissionsPolicy) Name() string {
	return "job_permissions"
}

func (policy *JobPermissionsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
	failed := false
	if wf.Permissions != nil && len(wf.Permissions) == 0 {
		return nil
	}
	if len(wf.Jobs) < 2 && wf.Permissions != nil {
		return nil
	}
	for jobName, job := range wf.Jobs {
		if job.Permissions == nil {
			failed = true
			logE.WithField("job_name", jobName).Error("job should have permissions")
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
