package policy

import (
	"errors"

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

func (p *JobPermissionsPolicy) ApplyJob(_ *logrus.Entry, _ *config.Config, jobCtx *JobContext, job *workflow.Job) error {
	wf := jobCtx.Workflow.Workflow
	wfPermissions := wf.Permissions.Permissions()
	if wfPermissions != nil && len(wfPermissions) == 0 {
		// workflow's permissions is `{}`
		return nil
	}
	if len(wf.Jobs) < 2 && wfPermissions != nil {
		// workflow permissions is set and there is only one job
		return nil
	}
	if job.Permissions.IsNil() {
		return errors.New("job should have permissions")
	}
	return nil
}
