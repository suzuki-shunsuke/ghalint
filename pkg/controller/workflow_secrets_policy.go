package controller

import (
	"context"
	"errors"
	"regexp"

	"github.com/sirupsen/logrus"
)

type WorkflowSecretsPolicy struct {
	secretPattern      *regexp.Regexp
	githubTokenPattern *regexp.Regexp
}

func NewWorkflowSecretsPolicy() *WorkflowSecretsPolicy {
	return &WorkflowSecretsPolicy{
		secretPattern:      regexp.MustCompile(`\${{ *secrets\.[^ ]+ *}}`),
		githubTokenPattern: regexp.MustCompile(`\${{ *github\.token+ *}}`),
	}
}

func (p *WorkflowSecretsPolicy) Name() string {
	return "workflow_secrets"
}

func (p *WorkflowSecretsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
	if len(wf.Jobs) < 2 { //nolint:gomnd
		return nil
	}
	failed := false
	for envName, envValue := range wf.Env {
		if p.secretPattern.MatchString(envValue) {
			failed = true
			logE.WithField("env_name", envName).Error("secret should not be set to workflow's env")
		}
		if p.githubTokenPattern.MatchString(envValue) {
			failed = true
			logE.WithField("env_name", envName).Error("github.token should not be set to workflow's env")
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
