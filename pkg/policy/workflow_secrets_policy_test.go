package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestWorkflowSecretsPolicy_ApplyWorkflow(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    *workflow.Workflow
		isErr bool
	}{
		{
			name: "workflow has only one job",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Jobs: map[string]*workflow.Job{
					"foo": {},
				},
			},
		},
		{
			name: "secret should not be set to workflow's env",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
				},
				Jobs: map[string]*workflow.Job{
					"foo": {},
					"bar": {},
				},
			},
			isErr: true,
		},
		{
			name: "github token should not be set to workflow's env",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Jobs: map[string]*workflow.Job{
					"foo": {},
					"bar": {},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Env: map[string]string{
					"FOO": "foo",
				},
				Jobs: map[string]*workflow.Job{
					"foo": {},
					"bar": {},
				},
			},
		},
	}
	p := policy.NewWorkflowSecretsPolicy()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyWorkflow(logE, d.cfg, nil, d.wf); err != nil {
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
