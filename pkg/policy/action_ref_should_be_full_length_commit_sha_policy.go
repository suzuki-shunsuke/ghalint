package policy

import (
	"context"
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

func (p *ActionRefShouldBeSHA1Policy) Apply(_ context.Context, logE *logrus.Entry, cfg *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if p.applyJob(logE, cfg, job) {
			failed = true
		}
	}
	if failed {
		return errors.New("workflow violates policies")
	}
	return nil
}

func (p *ActionRefShouldBeSHA1Policy) applyJob(logE *logrus.Entry, cfg *config.Config, job *workflow.Job) bool {
	if action := p.checkUses(job.Uses); action != "" {
		if p.excluded(action, cfg.Excludes) {
			return false
		}
		logE.WithField("uses", job.Uses).Error("action ref should be full length SHA1")
		return true
	}
	for _, step := range job.Steps {
		action := p.checkUses(step.Uses)
		if action == "" || p.excluded(action, cfg.Excludes) {
			continue
		}
		fields := logrus.Fields{
			"uses": step.Uses,
		}
		if step.ID != "" {
			fields["step_id"] = step.ID
		}
		if step.Name != "" {
			fields["step_name"] = step.Name
		}
		logE.WithFields(fields).Error("action ref should be full length SHA1")
		return true
	}
	return false
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
