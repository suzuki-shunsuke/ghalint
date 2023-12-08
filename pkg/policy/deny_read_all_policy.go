package policy

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type DenyReadAllPermissionPolicy struct{}

func (p *DenyReadAllPermissionPolicy) Name() string {
	return "deny_read_all_permission"
}

func (p *DenyReadAllPermissionPolicy) ID() string {
	return "002"
}

func (p *DenyReadAllPermissionPolicy) Apply(_ context.Context, logE *logrus.Entry, _ *config.Config, wf *workflow.Workflow) error {
	failed := false
	wfReadAll := wf.Permissions.ReadAll()
	for jobName, job := range wf.Jobs {
		if job.Permissions.ReadAll() {
			failed = true
			logE.WithField("job_name", jobName).Error("don't use read-all permission")
			continue
		}
		if job.Permissions.IsNil() && wfReadAll {
			failed = true
			logE.WithField("job_name", jobName).Error("don't use read-all permission")
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
