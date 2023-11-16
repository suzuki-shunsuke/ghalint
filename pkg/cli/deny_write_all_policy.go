package cli

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type DenyWriteAllPermissionPolicy struct{}

func (p *DenyWriteAllPermissionPolicy) Name() string {
	return "deny_write_all_permission"
}

func (p *DenyWriteAllPermissionPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
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
		return errors.New("workflow violates policies")
	}
	return nil
}
