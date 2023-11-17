package policy_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestActionRefShouldBeSHA1Policy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    *workflow.Workflow
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
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"release": {
						Steps: []*workflow.Step{
							{
								Uses: "slsa-framework/slsa-github-generator@v1.5.0",
							},
						},
					},
					"release2": {
						Uses: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@v0.4.4",
					},
				},
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
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"release": {
						Steps: []*workflow.Step{
							{
								Uses: "slsa-framework/slsa-github-generator@v1.5.0",
								ID:   "generate",
								Name: "Generate SLSA Provenance",
							},
						},
					},
				},
			},
		},
		{
			name:  "job error",
			isErr: true,
			cfg:   &config.Config{},
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"release": {
						Uses: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@v0.4.4",
					},
				},
			},
		},
	}
	p := policy.NewActionRefShouldBeSHA1Policy()
	ctx := context.Background()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
