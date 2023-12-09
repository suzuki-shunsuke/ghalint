package policy

import (
	"errors"

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

func (p *DenyWriteAllPermissionPolicy) ApplyJob(_ *logrus.Entry, _ *config.Config, jobCtx *JobContext, job *workflow.Job) error {
	wfWriteAll := jobCtx.Workflow.Workflow.Permissions.WriteAll()
	if job.Permissions.WriteAll() {
		return errors.New("don't use write-all permission")
	}
	if job.Permissions.IsNil() && wfWriteAll {
		return errors.New("don't use write-all permission")
	}
	return nil
}
