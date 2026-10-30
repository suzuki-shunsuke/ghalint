package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestActionShouldBeLockedPolicy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	lf := &lockfile.Lockfile{
		Workflows: map[string][]string{
			".github/workflows/test.yaml": {
				"actions/checkout@v6",
			},
		},
		Dependencies: map[string]*lockfile.Dependency{
			"actions/checkout@v6": {
				Ref:    "v6",
				Commit: "1af3b93b6815bc44a9784bd300feb67ff0d1eeb3",
			},
		},
	}
	stepCtx := func(lf *lockfile.Lockfile) *policy.StepContext {
		return &policy.StepContext{
			FilePath: ".github/workflows/test.yaml",
			Job: &policy.JobContext{
				Name: "test",
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
					Lockfile: lf,
				},
			},
		}
	}
	data := []struct {
		name    string
		cfg     *config.Config
		stepCtx *policy.StepContext
		step    *workflow.Step
		isErr   bool
	}{
		{
			name:    "locked",
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "actions/checkout@v6"},
		},
		{
			name:    "not locked",
			isErr:   true,
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "actions/setup-go@v6"},
		},
		{
			name:    "pinned by commit SHA",
			isErr:   true,
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "actions/checkout@1af3b93b6815bc44a9784bd300feb67ff0d1eeb3"},
		},
		{
			name:    "local action",
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "./.github/actions/foo"},
		},
		{
			name:    "docker image",
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "docker://rhysd/actionlint:latest"},
		},
		{
			name:    "run step",
			cfg:     &config.Config{},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Run: "echo hello"},
		},
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_should_be_locked",
						ActionName: "actions/*",
					},
				},
			},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "actions/setup-go@v6"},
		},
		{
			name: "lockfile is disabled",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "lockfile_is_required",
					},
				},
			},
			stepCtx: stepCtx(lf),
			step:    &workflow.Step{Uses: "actions/setup-go@v6"},
		},
		{
			name:    "lockfile doesn't exist",
			cfg:     &config.Config{},
			stepCtx: stepCtx(nil),
			step:    &workflow.Step{Uses: "actions/setup-go@v6"},
		},
		{
			name:    "action file",
			cfg:     &config.Config{},
			stepCtx: &policy.StepContext{Action: &workflow.Action{}},
			step:    &workflow.Step{Uses: "actions/setup-go@v6"},
		},
	}
	p := policy.NewActionShouldBeLockedPolicy()
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyStep(logger, d.cfg, d.stepCtx, d.step); err != nil {
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
