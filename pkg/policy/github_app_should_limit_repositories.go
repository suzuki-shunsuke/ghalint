package policy

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type GitHubAppShouldLimitRepositoriesPolicy struct{}

func (p *GitHubAppShouldLimitRepositoriesPolicy) Name() string {
	return "github_app_should_limit_repositories"
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) ID() string {
	return "009"
}

func (p *GitHubAppShouldLimitRepositoriesPolicy) ApplyStep(logE *logrus.Entry, cfg *config.Config, stepCtx *StepContext, step *workflow.Step) (ge error) { //nolint:cyclop
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
	if p.excluded(cfg, stepCtx, step) {
		logE.Debug("this step is ignored")
		return nil
	}
	if action == "tibdex/github-app-token" {
		if step.With == nil {
			return errRepositoriesIsRequired
		}
		if _, ok := step.With["repositories"]; !ok {
			return errRepositoriesIsRequired
		}
		return nil
	}
	if action == "actions/create-github-app-token" {
		if step.With == nil {
			return errRepositoriesIsRequired
		}
		if _, ok := step.With["repositories"]; ok {
			return nil
		}
		if _, ok := step.With["owner"]; ok {
			return errRepositoriesIsRequired
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

func (p *GitHubAppShouldLimitRepositoriesPolicy) excluded(cfg *config.Config, stepCtx *StepContext, step *workflow.Step) bool {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if exclude.FilePath() != stepCtx.FilePath {
			continue
		}
		if stepCtx.Job != nil && exclude.JobName != stepCtx.Job.Name {
			continue
		}
		if exclude.StepID != step.ID {
			continue
		}
		return true
	}
	return false
}
