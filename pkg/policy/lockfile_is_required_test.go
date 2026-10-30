package policy_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
)

func TestLockfileIsRequiredPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		cfg      *config.Config
		lockfile *lockfile.Lockfile
		isErr    bool
	}{
		{
			name:     "lockfile exists",
			cfg:      &config.Config{},
			lockfile: &lockfile.Lockfile{},
		},
		{
			name:  "lockfile doesn't exist",
			cfg:   &config.Config{},
			isErr: true,
		},
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "lockfile_is_required",
					},
				},
			},
		},
	}
	p := &policy.LockfileIsRequiredPolicy{}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(d.cfg, d.lockfile); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
