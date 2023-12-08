package policy

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type GitHubAppShouldLimitPermissionsPolicy struct{}

func (p *GitHubAppShouldLimitPermissionsPolicy) Name() string {
	return "github_app_should_limit_permissions"
}

func (p *GitHubAppShouldLimitPermissionsPolicy) ID() string {
	return "010"
}

func (p *GitHubAppShouldLimitPermissionsPolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if err := p.applyJob(logE, cfg, wf.FilePath, jobName, job); err != nil {
			failed = true
		}
	}
	if failed {
		return errWorkflowViolatePolicy
	}
	return nil
}

func (p *GitHubAppShouldLimitPermissionsPolicy) applyJob(logE *logrus.Entry, cfg *config.Config, wfFilePath, jobName string, job *workflow.Job) error {
	failed := false
	for _, step := range job.Steps {
		if err := p.applyStep(logE, cfg, wfFilePath, jobName, step); err != nil {
			logerr.WithError(logE, err).WithField("step_id", step.ID).Error(`the step violates the policy`)
			failed = true
		}
	}
	if failed {
		return errJobViolatePolicy
	}
	return nil
}

func (p *GitHubAppShouldLimitPermissionsPolicy) applyStep(logE *logrus.Entry, cfg *config.Config, wfFilePath, jobName string, step *workflow.Step) (ge error) {
	action := p.checkUses(step.Uses)
	if action == "" {
		return nil
	}
	defer func() {
		if ge != nil {
			ge = logerr.WithFields(ge, logrus.Fields{
				"action": action,
			})
		}
	}()
	if p.excluded(cfg.Excludes, wfFilePath, jobName, step.ID) {
		logE.Debug("this step is ignored")
		return nil
	}
	if action == "tibdex/github-app-token" {
		if step.With == nil {
			return errPermissionsIsRequired
		}
		if _, ok := step.With["permissions"]; !ok {
			return errPermissionsIsRequired
		}
		return nil
	}
	return nil
}

func (p *GitHubAppShouldLimitPermissionsPolicy) checkUses(uses string) string {
	if uses == "" {
		return ""
	}
	action, _, _ := strings.Cut(uses, "@")
	return action
}

func (p *GitHubAppShouldLimitPermissionsPolicy) excluded(excludes []*config.Exclude, wfFilePath, jobName, stepID string) bool {
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
