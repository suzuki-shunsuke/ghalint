package policy

import (
	"errors"
	"log/slog"
	"path"
	"regexp"
	"strings"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type ActionRefShouldBeSHAPolicy struct {
	sha1Pattern   *regexp.Regexp
	sha256Pattern *regexp.Regexp
}

func NewActionRefShouldBeSHAPolicy() *ActionRefShouldBeSHAPolicy {
	return &ActionRefShouldBeSHAPolicy{
		sha1Pattern:   regexp.MustCompile(`\b[0-9a-f]{40}\b`),
		sha256Pattern: regexp.MustCompile(`\b[0-9a-f]{64}\b`),
	}
}

func (p *ActionRefShouldBeSHAPolicy) Name() string {
	return "action_ref_should_be_full_length_commit_sha"
}

func (p *ActionRefShouldBeSHAPolicy) ID() string {
	return "008"
}

func (p *ActionRefShouldBeSHAPolicy) ApplyJob(_ *slog.Logger, cfg *config.Config, _ *JobContext, job *workflow.Job) error {
	return p.apply(cfg, job.Uses)
}

// ApplyStep validates steps of workflows and actions.
// Actions in workflows locked by the lockfile are ignored because the lockfile locks them to commits.
// If ghalint fails to read the lockfile, actions in workflows are ignored too
// because they may be locked and the read error already fails the run.
// Steps of action files are always validated because callers of the action files don't necessarily use the lockfile.
func (p *ActionRefShouldBeSHAPolicy) ApplyStep(_ *slog.Logger, cfg *config.Config, stepCtx *StepContext, step *workflow.Step) error {
	if p.locked(stepCtx, step.Uses) {
		return nil
	}
	return p.apply(cfg, step.Uses)
}

func (p *ActionRefShouldBeSHAPolicy) locked(stepCtx *StepContext, uses string) bool {
	if stepCtx == nil || stepCtx.Job == nil || stepCtx.Job.Workflow == nil {
		return false
	}
	wfCtx := stepCtx.Job.Workflow
	return wfCtx.Lockfile.Locked(wfCtx.FilePath, uses)
}

func (p *ActionRefShouldBeSHAPolicy) apply(cfg *config.Config, uses string) error {
	action := p.checkUses(uses)
	if action == "" || p.excluded(action, cfg.Excludes) {
		return nil
	}
	return slogerr.With(errors.New("action ref should be full length SHA"), //nolint:wrapcheck
		"action", action,
	)
}

func (p *ActionRefShouldBeSHAPolicy) checkUses(uses string) string {
	if uses == "" {
		return ""
	}
	if ref, ok := strings.CutPrefix(uses, "docker://"); ok {
		repoAndTag, digest, hasDigest := strings.Cut(ref, "@sha256:")
		if hasDigest && p.sha256Pattern.MatchString(digest) {
			return ""
		}
		repo := repoAndTag
		lastColon := strings.LastIndex(repoAndTag, ":")
		lastSlash := strings.LastIndex(repoAndTag, "/")
		if lastColon != -1 && lastColon > lastSlash {
			repo = repoAndTag[:lastColon]
		}
		return "docker://" + repo
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
