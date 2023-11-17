package controller_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
)

func TestWorkflowSecretsPolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *controller.Config
		wf    *controller.Workflow
		isErr bool
	}{
		{
			name: "workflow has only one job",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Jobs: map[string]*controller.Job{
					"foo": {},
				},
			},
		},
		{
			name: "secret should not be set to workflow's env",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
				},
				Jobs: map[string]*controller.Job{
					"foo": {},
					"bar": {},
				},
			},
			isErr: true,
		},
		{
			name: "github token should not be set to workflow's env",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Jobs: map[string]*controller.Job{
					"foo": {},
					"bar": {},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"FOO": "foo",
				},
				Jobs: map[string]*controller.Job{
					"foo": {},
					"bar": {},
				},
			},
		},
	}
	policy := controller.NewWorkflowSecretsPolicy()
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
