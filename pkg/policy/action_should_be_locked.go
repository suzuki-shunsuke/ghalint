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

// ActionShouldBeLockedPolicy requires actions in workflows to be managed by the lockfile.
// Actions pinned by full length commit SHA aren't allowed because the lockfile locks actions recursively.
type ActionShouldBeLockedPolicy struct {
	sha1Pattern *regexp.Regexp
}

func NewActionShouldBeLockedPolicy() *ActionShouldBeLockedPolicy {
	return &ActionShouldBeLockedPolicy{
		sha1Pattern: regexp.MustCompile(`^[0-9a-f]{40}$`),
	}
}

func (p *ActionShouldBeLockedPolicy) Name() string {
	return "action_should_be_locked"
}

func (p *ActionShouldBeLockedPolicy) ID() string {
	return "015"
}

func (p *ActionShouldBeLockedPolicy) ApplyStep(_ *slog.Logger, cfg *config.Config, stepCtx *StepContext, step *workflow.Step) error {
	if stepCtx.Job == nil || !LockfileEnabled(cfg) {
		// Actions in action files aren't validated because callers of the action files don't necessarily use the lockfile.
		// They should be pinned by full length commit SHA, which action_ref_should_be_full_length_commit_sha validates.
		return nil
	}
	lf := stepCtx.Job.Workflow.Lockfile
	if lf == nil {
		// lockfile_is_required reports the missing lockfile.
		return nil
	}
	action, ref, ok := p.parseUses(step.Uses)
	if !ok || p.excluded(action, cfg.Excludes) {
		return nil
	}
	if p.sha1Pattern.MatchString(ref) {
		return slogerr.With(errors.New("action shouldn't be pinned by commit SHA but be managed by the lockfile"), //nolint:wrapcheck
			"action", action,
		)
	}
	if lf.Locked(stepCtx.Job.Workflow.FilePath, step.Uses) {
		return nil
	}
	return slogerr.With(errors.New("action should be managed by the lockfile"), //nolint:wrapcheck
		"action", step.Uses,
	)
}

// parseUses returns the action name and ref of a remote action.
// It returns false for local actions and Docker images, which the lockfile doesn't manage.
func (p *ActionShouldBeLockedPolicy) parseUses(uses string) (string, string, bool) {
	if uses == "" || strings.HasPrefix(uses, "./") || strings.HasPrefix(uses, "docker://") {
		return "", "", false
	}
	return strings.Cut(uses, "@")
}

func (p *ActionShouldBeLockedPolicy) excluded(action string, excludes []*config.Exclude) bool {
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
