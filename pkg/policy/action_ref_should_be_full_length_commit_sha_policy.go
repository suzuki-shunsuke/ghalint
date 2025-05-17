package policy

import (
	"errors"
	"path"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type ActionRefShouldBeSHAPolicy struct {
	sha1Pattern *regexp.Regexp
}

func NewActionRefShouldBeSHAPolicy() *ActionRefShouldBeSHAPolicy {
	return &ActionRefShouldBeSHAPolicy{
		sha1Pattern: regexp.MustCompile(`\b[0-9a-f]{40}\b`),
	}
}

func (p *ActionRefShouldBeSHAPolicy) Name() string {
	return "action_ref_should_be_full_length_commit_sha"
}

func (p *ActionRefShouldBeSHAPolicy) ID() string {
	return "008"
}

func (p *ActionRefShouldBeSHAPolicy) ApplyJob(_ *logrus.Entry, cfg *config.Config, _ *JobContext, job *workflow.Job) error {
	return p.apply(cfg, job.Uses)
}

func (p *ActionRefShouldBeSHAPolicy) ApplyStep(_ *logrus.Entry, cfg *config.Config, _ *StepContext, step *workflow.Step) error {
	return p.apply(cfg, step.Uses)
}

func (p *ActionRefShouldBeSHAPolicy) apply(cfg *config.Config, uses string) error {
	action := p.checkUses(uses)
	if action == "" || p.excluded(action, cfg.Excludes) {
		return nil
	}
	return logerr.WithFields(errors.New("action ref should be full length SHA1"), logrus.Fields{ //nolint:wrapcheck
		"action": action,
	})
}

func (p *ActionRefShouldBeSHAPolicy) checkUses(uses string) string {
	if uses == "" {
		return ""
	}
	action, tag, ok := strings.Cut(uses, "@")
	if !ok {
		return ""
	}
	if p.sha1Pattern.MatchString(tag) {
		return ""
	}
	return action
}

func (p *ActionRefShouldBeSHAPolicy) excluded(action string, excludes []*config.Exclude) bool {
	for _, exclude := range excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if f, _ := path.Match(exclude.ActionName, action); f {
			return true
		}
	}
	return false
}
