package policy

import (
	"context"
	"errors"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type JobSecretsPolicy struct {
	secretPattern      *regexp.Regexp
	githubTokenPattern *regexp.Regexp
}

func NewJobSecretsPolicy() *JobSecretsPolicy {
	return &JobSecretsPolicy{
		secretPattern:      regexp.MustCompile(`\${{ *secrets\.[^ ]+ *}}`),
		githubTokenPattern: regexp.MustCompile(`\${{ *github\.token+ *}}`),
	}
}

func (p *JobSecretsPolicy) Name() string {
	return "job_secrets"
}

func (p *JobSecretsPolicy) ID() string {
	return "006"
}

func checkExcludes(policyName string, wf *workflow.Workflow, jobName string, cfg *config.Config) bool {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == policyName && wf.FilePath == exclude.WorkflowFilePath && jobName == exclude.JobName {
			return true
		}
	}
	return false
}

func (p *JobSecretsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if checkExcludes(p.Name(), wf, jobName, cfg) {
			continue
		}
		if len(job.Steps) < 2 { //nolint:gomnd
			continue
		}
		for envName, envValue := range job.Env {
			if p.secretPattern.MatchString(envValue) {
				failed = true
				logE.WithField("env_name", envName).Error("secret should not be set to job's env")
			}
			if p.githubTokenPattern.MatchString(envValue) {
				failed = true
				logE.WithField("env_name", envName).Error("github.token should not be set to job's env")
			}
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}
