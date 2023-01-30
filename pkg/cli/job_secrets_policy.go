package cli

import (
	"context"
	"errors"
	"regexp"

	"github.com/sirupsen/logrus"
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

func (policy *JobSecretsPolicy) Name() string {
	return "job_secrets"
}

func checkExcludes(policyName string, wf *Workflow, jobName string, cfg *Config) bool {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == policyName && wf.FilePath == exclude.WorkflowFilePath && jobName == exclude.JobName {
			return true
		}
	}
	return false
}

func (policy *JobSecretsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
	if len(wf.Jobs) < 2 { //nolint:gomnd
		return nil
	}
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if checkExcludes(policy.Name(), wf, jobName, cfg) {
			continue
		}
		if len(job.Steps) < 2 { //nolint:gomnd
			continue
		}
		for envName, envValue := range job.Env {
			if policy.secretPattern.MatchString(envValue) {
				failed = true
				logE.WithField("env_name", envName).Error("secret should not be set to job's env")
			}
			if policy.githubTokenPattern.MatchString(envValue) {
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
