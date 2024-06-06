package policy

import (
	"errors"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
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

func checkExcludes(policyName string, jobCtx *JobContext, cfg *config.Config) bool {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == policyName && jobCtx.Workflow.FilePath == exclude.WorkflowFilePath && jobCtx.Name == exclude.JobName {
			return true
		}
	}
	return false
}

func (p *JobSecretsPolicy) ApplyJob(_ *logrus.Entry, cfg *config.Config, jobCtx *JobContext, job *workflow.Job) error {
	if checkExcludes(p.Name(), jobCtx, cfg) {
		return nil
	}
	if len(job.Steps) < 2 { //nolint:mnd
		return nil
	}
	for envName, envValue := range job.Env {
		if p.secretPattern.MatchString(envValue) {
			return logerr.WithFields(errors.New("secret should not be set to job's env"), logrus.Fields{ //nolint:wrapcheck
				"env_name": envName,
			})
		}
		if p.githubTokenPattern.MatchString(envValue) {
			return logerr.WithFields(errors.New("github.token should not be set to job's env"), logrus.Fields{ //nolint:wrapcheck
				"env_name": envName,
			})
		}
	}
	return nil
}
