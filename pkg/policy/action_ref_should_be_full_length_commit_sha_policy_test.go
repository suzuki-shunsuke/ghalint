package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestActionRefShouldBeSHAPolicy_ApplyJob(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		job   *workflow.Job
		isErr bool
	}{
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "slsa-framework/slsa-github-generator",
					},
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml",
					},
				},
			},
			job: &workflow.Job{
				Steps: []*workflow.Step{
					{
						Uses: "slsa-framework/slsa-github-generator@v1.5.0",
					},
				},
			},
		},
		{
			name:  "job error",
			isErr: true,
			cfg:   &config.Config{},
			job: &workflow.Job{
				Uses: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@v0.4.4",
			},
		},
		{
			name: "docker image with digest",
			cfg:  &config.Config{},
			job: &workflow.Job{
				Uses: "docker://rhysd/actionlint:1.7.7@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name: "docker image with digest (no tag)",
			cfg:  &config.Config{},
			job: &workflow.Job{
				Uses: "docker://rhysd/actionlint@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name: "docker image with port and digest",
			cfg:  &config.Config{},
			job: &workflow.Job{
				Uses: "docker://registry.example.com:5000/myimage:1.0.0@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name:  "docker image with tag",
			isErr: true,
			cfg:   &config.Config{},
			job: &workflow.Job{
				Uses: "docker://rhysd/actionlint:latest",
			},
		},
		{
			name:  "docker image with port and tag",
			isErr: true,
			cfg:   &config.Config{},
			job: &workflow.Job{
				Uses: "docker://registry.example.com:5000/myimage:latest",
			},
		},
		{
			name: "exclude docker image with tag",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "docker://rhysd/actionlint",
					},
				},
			},
			job: &workflow.Job{
				Uses: "docker://rhysd/actionlint:latest",
			},
		},
		{
			name: "exclude docker image with port and tag",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "docker://registry.example.com:5000/myimage",
					},
				},
			},
			job: &workflow.Job{
				Uses: "docker://registry.example.com:5000/myimage:latest",
			},
		},
	}
	p := policy.NewActionRefShouldBeSHAPolicy()
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logger, d.cfg, nil, d.job); err != nil {
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

func TestActionRefShouldBeSHAPolicy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	lf := testLockfile(t)
	data := []struct {
		name    string
		cfg     *config.Config
		stepCtx *policy.StepContext
		step    *workflow.Step
		isErr   bool
	}{
		{
			name:    "locked by the lockfile",
			cfg:     &config.Config{},
			stepCtx: workflowStepContext(lf),
			step: &workflow.Step{
				Uses: "actions/checkout@v6",
			},
		},
		{
			name:    "not locked by the lockfile",
			isErr:   true,
			cfg:     &config.Config{},
			stepCtx: workflowStepContext(lf),
			step: &workflow.Step{
				Uses: "actions/setup-go@v6",
			},
		},
		{
			name:    "lockfile doesn't exist",
			isErr:   true,
			cfg:     &config.Config{},
			stepCtx: workflowStepContext(nil),
			step: &workflow.Step{
				Uses: "actions/checkout@v6",
			},
		},
		{
			name: "lockfile is invalid",
			cfg:  &config.Config{},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Workflow: &policy.WorkflowContext{
						FilePath:        ".github/workflows/test.yaml",
						LockfileInvalid: true,
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/setup-go@v6",
			},
		},
		{
			name:  "docker image is validated even if lockfile is invalid",
			isErr: true,
			cfg:   &config.Config{},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Workflow: &policy.WorkflowContext{
						FilePath:        ".github/workflows/test.yaml",
						LockfileInvalid: true,
					},
				},
			},
			step: &workflow.Step{
				Uses: "docker://rhysd/actionlint:latest",
			},
		},
		{
			name:    "action file",
			isErr:   true,
			cfg:     &config.Config{},
			stepCtx: &policy.StepContext{Action: &workflow.Action{}},
			step: &workflow.Step{
				Uses: "actions/checkout@v6",
			},
		},
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "slsa-framework/slsa-github-generator",
					},
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "slsa-framework/slsa-github-generator@v1.5.0",
			},
		},
		{
			name: "exclude with glob pattern",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "slsa-framework/*",
					},
				},
			},
			step: &workflow.Step{
				Uses: "slsa-framework/slsa-github-generator@v1.5.0",
			},
		},
		{
			name:  "step error",
			isErr: true,
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "actions/checkout",
					},
				},
			},
			step: &workflow.Step{
				Uses: "slsa-framework/slsa-github-generator@v1.5.0",
				ID:   "generate",
				Name: "Generate SLSA Provenance",
			},
		},
		{
			name: "docker image with digest",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "docker://rhysd/actionlint:1.7.7@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name: "docker image with digest (no tag)",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "docker://rhysd/actionlint@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name: "docker image with port and digest",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "docker://registry.example.com:5000/myimage:1.0.0@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9",
			},
		},
		{
			name:  "docker image with tag",
			isErr: true,
			cfg:   &config.Config{},
			step: &workflow.Step{
				Uses: "docker://rhysd/actionlint:latest",
			},
		},
		{
			name:  "docker image with port and tag",
			isErr: true,
			cfg:   &config.Config{},
			step: &workflow.Step{
				Uses: "docker://registry.example.com:5000/myimage:latest",
			},
		},
		{
			name: "exclude docker image with tag",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "docker://rhysd/actionlint",
					},
				},
			},
			step: &workflow.Step{
				Uses: "docker://rhysd/actionlint:latest",
			},
		},
		{
			name: "exclude docker image with port and tag",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "docker://registry.example.com:5000/myimage",
					},
				},
			},
			step: &workflow.Step{
				Uses: "docker://registry.example.com:5000/myimage:latest",
			},
		},
	}
	p := policy.NewActionRefShouldBeSHAPolicy()
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

func testLockfile(t *testing.T) *lockfile.Lockfile {
	t.Helper()
	lf, err := lockfile.Parse([]byte(`version: 'v0.0.3'
workflows:
  '.github/workflows/test.yaml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    commit: 'sha1-3d3c42e5aac5ba805825da76410c181273ba90b1'
    owner_id: 44036562
    repo_id: 197814629
`))
	if err != nil {
		t.Fatal(err)
	}
	return lf
}

func workflowStepContext(lf *lockfile.Lockfile) *policy.StepContext {
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
