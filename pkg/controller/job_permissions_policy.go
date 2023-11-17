package controller

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type JobPermissionsPolicy struct{}

func (p *JobPermissionsPolicy) Name() string {
	return "job_permissions"
}

func (p *JobPermissionsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
	failed := false
	wfPermissions := wf.Permissions.Permissions()
	if wfPermissions != nil && len(wfPermissions) == 0 {
		// workflow's permissions is `{}`
		return nil
	}
	if len(wf.Jobs) < 2 && wfPermissions != nil {
		// workflow permissions is set and there is only one job
		return nil
	}
	for jobName, job := range wf.Jobs {
		if job.Permissions.IsNil() {
			failed = true
			logE.WithField("job_name", jobName).Error("job should have permissions")
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
