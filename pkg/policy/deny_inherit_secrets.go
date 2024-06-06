package policy

import (
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

func (p *DenyInheritSecretsPolicy) ApplyJob(_ *logrus.Entry, cfg *config.Config, jobCtx *JobContext, job *workflow.Job) error {
	if checkExcludes(p.Name(), jobCtx, cfg) {
		return nil
	}
	if job.Secrets.Inherit() {
		return errors.New("`secrets: inherit` should not be used. Only required secrets should be passed explicitly")
	}
	return nil
}
