package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestActionRefShouldBeSHA1Policy_ApplyJob(t *testing.T) {
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
	}
	p := policy.NewActionRefShouldBeSHA1Policy()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logE, d.cfg, nil, d.job); err != nil {
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

func TestActionRefShouldBeSHA1Policy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		step  *workflow.Step
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
	}
	p := policy.NewActionRefShouldBeSHA1Policy()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyStep(logE, d.cfg, nil, d.step); err != nil {
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
