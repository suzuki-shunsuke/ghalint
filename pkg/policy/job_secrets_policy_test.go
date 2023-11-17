package policy_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestJobSecretsPolicy_Apply(t *testing.T) { //nolint:funlen
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
						PolicyName:       "job_secrets",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "foo",
					},
				},
			},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*workflow.Step{
							{},
							{},
						},
					},
				},
			},
		},
		{
			name: "job has only one step",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*workflow.Step{
							{},
						},
					},
				},
			},
		},
		{
			name: "secret should not be set to job's env",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
						},
						Steps: []*workflow.Step{
							{},
							{},
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "github token should not be set to job's env",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*workflow.Step{
							{},
							{},
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"foo": {
						Env: map[string]string{
							"FOO": "foo",
						},
						Steps: []*workflow.Step{
							{},
							{},
						},
					},
				},
			},
		},
	}
	p := policy.NewJobSecretsPolicy()
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
