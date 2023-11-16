package cli_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
)

func TestActionRefShouldBeSHA1Policy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *cli.Config
		wf    *cli.Workflow
		isErr bool
	}{
		{
			name: "exclude",
			cfg: &cli.Config{
				Excludes: []*cli.Exclude{
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
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"release": {
						Steps: []*cli.Step{
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
			cfg: &cli.Config{
				Excludes: []*cli.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
						ActionName: "actions/checkout",
					},
				},
			},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"release": {
						Steps: []*cli.Step{
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
			cfg:   &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"release": {
						Uses: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@v0.4.4",
					},
				},
			},
		},
	}
	p := cli.NewActionRefShouldBeSHA1Policy()
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
