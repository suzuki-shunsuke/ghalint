package policy

import (
	"context"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type GitHubAppShouldLimitRepositoriesPolicy struct{}

func (p *GitHubAppShouldLimitRepositoriesPolicy) Name() string {
	return "github_app_should_limit_repositories"
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) ID() string {
	return "009"
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if p.applyJob(logE, cfg, wf.FilePath, jobName, job) {
			failed = true
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) applyJob(logE *logrus.Entry, cfg *config.Config, wfFilePath, jobName string, job *workflow.Job) bool {
	failed := false
	for _, step := range job.Steps {
		if err := p.applyStep(logE, cfg, wfFilePath, jobName, step); err != nil {
			logE.WithError(err).Error(`the input "repositories" is required`)
			failed = true
		}
	}
	return failed
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) applyStep(logE *logrus.Entry, cfg *config.Config, wfFilePath, jobName string, step *workflow.Step) error {
	action := p.checkUses(step.Uses)
	if action == "" {
		return nil
	}
	if p.excluded(cfg.Excludes, wfFilePath, jobName, step.ID) {
		logE.Debug("this step is ignored")
		return nil
	}
	if action == "tibdex/github-app-token" {
		if step.With == nil {
			return errors.New(`the input "repositories" is required`)
		}
		if _, ok := step.With["repositories"]; !ok {
			return errors.New(`the input "repositories" is required`)
		}
		return nil
	}
	if action == "actions/create-github-app-token" {
		if step.With == nil {
			return errors.New(`the input "repositories" is required`)
		}
		if _, ok := step.With["repositories"]; ok {
			return nil
		}
		if _, ok := step.With["owner"]; ok {
			return errors.New(`the input "repositories" is required`)
		}
		return nil
	}
	return nil
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) checkUses(uses string) string {
	if uses == "" {
		return ""
	}
	action, _, _ := strings.Cut(uses, "@")
	return action
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) excluded(excludes []*config.Exclude, wfFilePath, jobName, stepID string) bool {
	for _, exclude := range excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if exclude.WorkflowFilePath != wfFilePath {
			continue
		}
		if exclude.JobName != jobName {
			continue
		}
		if exclude.StepID != stepID {
			continue
		}
		return true
	}
	return false
}
