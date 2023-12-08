package policy

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type JobPermissionsPolicy struct{}

func (p *JobPermissionsPolicy) Name() string {
	return "job_permissions"
}

func (p *JobPermissionsPolicy) ID() string {
	return "001"
}

func (p *JobPermissionsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
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
		return errWorkflowViolatePolicy
	}
	return nil
}
