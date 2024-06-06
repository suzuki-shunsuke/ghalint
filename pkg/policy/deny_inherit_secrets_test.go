//nolint:funlen
package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestDenyInheritSecretsPolicy_ApplyJob(t *testing.T) {
	t.Parallel()
	data := []struct {
		name   string
		job    string
		cfg    *config.Config
		jobCtx *policy.JobContext
		isErr  bool
	}{
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "deny_inherit_secrets",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "foo",
					},
				},
			},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "foo",
			},
			job: `secrets: inherit`,
		},
		{
			name: "not exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "deny_inherit_secrets",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "bar",
					},
				},
			},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "foo",
			},
			job:   `secrets: inherit`,
			isErr: true,
		},
		{
			name: "error",
			job:  `secrets: inherit`,
			cfg:  &config.Config{},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "foo",
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &config.Config{},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "foo",
			},
			job: `secrets:
      foo: ${{secrets.API_KEY}}`,
		},
	}
	p := &policy.DenyInheritSecretsPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			job := &workflow.Job{}
			if err := yaml.Unmarshal([]byte(d.job), job); err != nil {
				t.Fatal(err)
			}
			if err := p.ApplyJob(logE, d.cfg, d.jobCtx, job); err != nil {
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
