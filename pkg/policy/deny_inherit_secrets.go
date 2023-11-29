package policy

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type DenyInheritSecretsPolicy struct{}

func (p *DenyInheritSecretsPolicy) Name() string {
	return "deny_inherit_secrets"
}

func (p *DenyInheritSecretsPolicy) ID() string {
	return "004"
}

func (p *DenyInheritSecretsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if job.Secrets.Inherit() {
			failed = true
			logE.Error("`secrets: inherit` should not be used. Only required secrets should be passed explicitly")
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
