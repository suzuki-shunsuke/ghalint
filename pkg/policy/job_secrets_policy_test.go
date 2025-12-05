package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestJobSecretsPolicy_ApplyJob(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		cfg    *config.Config
		jobCtx *policy.JobContext
		job    *workflow.Job
		isErr  bool
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
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "foo",
			},
			job: &workflow.Job{
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Steps: []*workflow.Step{
					{},
					{},
				},
			},
		},
		{
			name:   "job has only one step",
			cfg:    &config.Config{},
			jobCtx: &policy.JobContext{},
			job: &workflow.Job{
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Steps: []*workflow.Step{
					{},
				},
			},
		},
		{
			name:   "secret should not be set to job's env",
			cfg:    &config.Config{},
			jobCtx: &policy.JobContext{},
			job: &workflow.Job{
				Env: map[string]string{
					"GITHUB_TOKEN": "${{secrets.GITHUB_TOKEN}}",
				},
				Steps: []*workflow.Step{
					{},
					{},
				},
			},
			isErr: true,
		},
		{
			name:   "github token should not be set to job's env",
			cfg:    &config.Config{},
			jobCtx: &policy.JobContext{},
			job: &workflow.Job{
				Env: map[string]string{
					"GITHUB_TOKEN": "${{github.token}}",
				},
				Steps: []*workflow.Step{
					{},
					{},
				},
			},
			isErr: true,
		},
		{
			name:   "pass",
			cfg:    &config.Config{},
			jobCtx: &policy.JobContext{},
			job: &workflow.Job{
				Env: map[string]string{
					"FOO": "foo",
				},
				Steps: []*workflow.Step{
					{},
					{},
				},
			},
		},
	}
	p := policy.NewJobSecretsPolicy()
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logger, d.cfg, d.jobCtx, d.job); err != nil {
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
