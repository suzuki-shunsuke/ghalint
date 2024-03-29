package policy

import (
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

func (p *GitHubAppShouldLimitPermissionsPolicy) ApplyStep(logE *logrus.Entry, cfg *config.Config, stepCtx *StepContext, step *workflow.Step) (ge error) {
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

	var name string
	if stepCtx.Job != nil {
		name = stepCtx.Job.Name
	}
	if p.excluded(cfg.Excludes, stepCtx.FilePath, name, step.ID) {
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

func (p *GitHubAppShouldLimitPermissionsPolicy) excluded(excludes []*config.Exclude, filePath, jobName, stepID string) bool {
	for _, exclude := range excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if exclude.WorkflowFilePath != filePath {
			continue
		}
		if jobName != "" {
			if exclude.JobName != jobName {
				continue
			}
		}
		if exclude.StepID != stepID {
			continue
		}
		return true
	}
	return false
}
