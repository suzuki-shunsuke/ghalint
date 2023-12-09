package policy

import (
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

func (p *DenyReadAllPermissionPolicy) ApplyJob(_ *logrus.Entry, _ *config.Config, jobCtx *JobContext, job *workflow.Job) error {
	wfReadAll := jobCtx.Workflow.Workflow.Permissions.ReadAll()
	if job.Permissions.ReadAll() {
		return errors.New("don't use read-all permission")
	}
	if job.Permissions.IsNil() && wfReadAll {
		return errors.New("don't use read-all permission")
	}
	return nil
}
