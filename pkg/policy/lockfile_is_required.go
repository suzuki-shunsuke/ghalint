package policy

import (
	"errors"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

// LockfileIsRequiredPolicy requires the workflow dependency lockfile.
// Excluding this policy disables the lockfile support, so actions are pinned by full length commit SHA instead.
type LockfileIsRequiredPolicy struct{}

const lockfileIsRequiredPolicyName = "lockfile_is_required"

func (p *LockfileIsRequiredPolicy) Name() string {
	return lockfileIsRequiredPolicyName
}

func (p *LockfileIsRequiredPolicy) ID() string {
	return "014"
}

// Apply returns an error if the lockfile doesn't exist.
// lf is nil if the lockfile doesn't exist.
func (p *LockfileIsRequiredPolicy) Apply(cfg *config.Config, lf *lockfile.Lockfile) error {
	if !LockfileEnabled(cfg) || lf != nil {
		return nil
	}
	return slogerr.With(errors.New("the lockfile is required"), //nolint:wrapcheck
		"lockfile", lockfile.Path,
	)
}

// LockfileEnabled reports whether actions in workflows should be managed by the lockfile.
// It's enabled unless the policy lockfile_is_required is excluded.
func LockfileEnabled(cfg *config.Config) bool {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == lockfileIsRequiredPolicyName {
			return false
		}
	}
	return true
}
