package cli_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
)

func TestJobSecretsPolicy_Apply(t *testing.T) { //nolint:funlen
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
						PolicyName:       "job_secrets",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "foo",
					},
				},
			},
			wf: &cli.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*cli.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []interface{}{
							map[string]interface{}{
								"run": "echo hello",
							},
							map[string]interface{}{
								"run": "echo bar",
							},
						},
					},
				},
			},
		},
		{
			name: "job has only one step",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*cli.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []interface{}{
							map[string]interface{}{
								"run": "echo hello",
							},
						},
					},
				},
			},
		},
		{
			name: "secret should not be set to job's env",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*cli.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
						},
						Steps: []interface{}{
							map[string]interface{}{
								"run": "echo hello",
							},
							map[string]interface{}{
								"run": "echo bar",
							},
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "github token should not be set to job's env",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*cli.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []interface{}{
							map[string]interface{}{
								"run": "echo hello",
							},
							map[string]interface{}{
								"run": "echo bar",
							},
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*cli.Job{
					"foo": {
						Env: map[string]string{
							"FOO": "foo",
						},
						Steps: []interface{}{
							map[string]interface{}{
								"run": "echo hello",
								"env": map[string]string{
									"GITHUB_TOKEN": "${{github.token}}",
								},
							},
							map[string]interface{}{
								"run": "echo bar",
							},
						},
					},
				},
			},
		},
	}
	policy := cli.NewJobSecretsPolicy()
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := policy.Apply(ctx, logE, d.cfg, d.wf); err != nil {
				if !d.isErr {
					t.Fatal(err)
				}
				return
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
