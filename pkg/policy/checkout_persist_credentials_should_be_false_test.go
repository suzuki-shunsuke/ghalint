package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestCheckoutPersistCredentialShouldBeFalsePolicy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		cfg     *config.Config
		step    *workflow.Step
		stepCtx *policy.StepContext
		isErr   bool
	}{
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "checkout_persist_credentials_should_be_false",
						WorkflowFilePath: ".github/workflows/test.yml",
						JobName:          "test",
					},
				},
			},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/checkout@v4",
			},
		},
		{
			name: "persist-credentials is not set",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "checkout_persist_credentials_should_be_false",
						JobName:          "test-2",
						WorkflowFilePath: ".github/workflows/test.yml",
					},
				},
			},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/checkout@v4",
			},
			isErr: true,
		},
		{
			name: "persist-credentials is true",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "checkout_persist_credentials_should_be_false",
						JobName:    "test-2",
					},
				},
			},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/checkout@v4",
				With: map[string]string{
					"persist-credentials": "true",
				},
			},
			isErr: true,
		},
		{
			name: "persist-credentials is false",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "checkout_persist_credentials_should_be_false",
						JobName:    "test-2",
					},
				},
			},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/checkout@v4",
				With: map[string]string{
					"persist-credentials": "false",
				},
			},
		},
		{
			name: "not checkout",
			cfg: &config.Config{
				Excludes: []*config.Exclude{},
			},
			stepCtx: &policy.StepContext{
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "actions/cache@v4",
			},
		},
	}
	p := &policy.CheckoutPersistCredentialShouldBeFalsePolicy{}
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
