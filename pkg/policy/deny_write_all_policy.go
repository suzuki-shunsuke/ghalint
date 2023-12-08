package policy

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type DenyWriteAllPermissionPolicy struct{}

func (p *DenyWriteAllPermissionPolicy) Name() string {
	return "deny_write_all_permission"
}

func (p *DenyWriteAllPermissionPolicy) ID() string {
	return "003"
}

func (p *DenyWriteAllPermissionPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	wfWriteAll := wf.Permissions.WriteAll()
	for jobName, job := range wf.Jobs {
		if job.Permissions.WriteAll() {
			failed = true
			logE.WithField("job_name", jobName).Error("don't use write-all permission")
			continue
		}
		if job.Permissions.IsNil() && wfWriteAll {
			failed = true
			logE.WithField("job_name", jobName).Error("don't use write-all permission")
		}
	}
	if failed {
		return workflowViolatePolicyError
	}
	return nil
}
