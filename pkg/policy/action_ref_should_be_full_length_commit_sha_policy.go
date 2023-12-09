package policy

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type ActionRefShouldBeSHA1Policy struct {
	sha1Pattern *regexp.Regexp
}

func NewActionRefShouldBeSHA1Policy() *ActionRefShouldBeSHA1Policy {
	return &ActionRefShouldBeSHA1Policy{
		sha1Pattern: regexp.MustCompile(`\b[0-9a-f]{40}\b`),
	}
}

func (p *ActionRefShouldBeSHA1Policy) Name() string {
	return "action_ref_should_be_full_length_commit_sha"
}

func (p *ActionRefShouldBeSHA1Policy) ID() string {
	return "008"
}

func (p *ActionRefShouldBeSHA1Policy) excluded(action string, excludes []*config.Exclude) bool {
	for _, exclude := range excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if action == exclude.ActionName {
			return true
		}
	}
	return false
}

func (p *ActionRefShouldBeSHA1Policy) ApplyJob(_ *logrus.Entry, cfg *config.Config, _ *JobContext, job *workflow.Job) error {
	action := p.checkUses(job.Uses)
	if action == "" {
		return nil
	}
	if p.excluded(action, cfg.Excludes) {
		return nil
	}
	return errors.New("action ref should be full length SHA1")
}

func (p *ActionRefShouldBeSHA1Policy) ApplyStep(_ *logrus.Entry, cfg *config.Config, _ *StepContext, step *workflow.Step) error {
	action := p.checkUses(step.Uses)
	if action == "" || p.excluded(action, cfg.Excludes) {
		return nil
	}
	return errors.New("action ref should be full length SHA1")
}

func (p *ActionRefShouldBeSHA1Policy) checkUses(uses string) string {
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
