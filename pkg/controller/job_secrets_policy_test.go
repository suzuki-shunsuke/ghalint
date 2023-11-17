package controller_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
)

func TestJobSecretsPolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *controller.Config
		wf    *controller.Workflow
		isErr bool
	}{
		{
			name: "exclude",
			cfg: &controller.Config{
				Excludes: []*controller.Exclude{
					{
						PolicyName:       "job_secrets",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "foo",
					},
				},
			},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*controller.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*controller.Step{
							{},
							{},
						},
					},
				},
			},
		},
		{
			name: "job has only one step",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*controller.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*controller.Step{
							{},
						},
					},
				},
			},
		},
		{
			name: "secret should not be set to job's env",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*controller.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
						},
						Steps: []*controller.Step{
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
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*controller.Job{
					"foo": {
						Env: map[string]string{
							"GITHUB_TOKEN": "${{github.token}}",
						},
						Steps: []*controller.Step{
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
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*controller.Job{
					"foo": {
						Env: map[string]string{
							"FOO": "foo",
						},
						Steps: []*controller.Step{
							{},
							{},
						},
					},
				},
			},
		},
	}
	policy := controller.NewJobSecretsPolicy()
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
