package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestJobPermissionsPolicy_ApplyJob(t *testing.T) { //nolint:funlen
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
						PolicyName:       "job_permissions",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "foo",
					},
				},
			},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
					Workflow: &workflow.Workflow{
						Permissions: &workflow.Permissions{},
						Jobs: map[string]*workflow.Job{
							"foo": {},
							"bar": {},
						},
					},
				},
				Name: "foo",
			},
			job: &workflow.Job{},
		},
		{
			name: "not exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "job_permissions",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "bar",
					},
				},
			},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
					Workflow: &workflow.Workflow{
						Permissions: &workflow.Permissions{},
						Jobs: map[string]*workflow.Job{
							"foo": {},
							"bar": {},
						},
					},
				},
				Name: "foo",
			},
			job:   &workflow.Job{},
			isErr: true,
		},
		{
			name: "workflow permissions is empty",
			job:  &workflow.Job{},
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{
						Permissions: workflow.NewPermissions(false, false, map[string]string{}),
						Jobs: map[string]*workflow.Job{
							"foo": {},
						},
					},
				},
			},
		},
		{
			name: "workflow has only one job",
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{
						Permissions: workflow.NewPermissions(false, false, map[string]string{
							"contents": "read",
						}),
						Jobs: map[string]*workflow.Job{
							"foo": {},
						},
					},
				},
			},
			job: &workflow.Job{},
		},
		{
			name: "job should have permissions",
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{
						Permissions: &workflow.Permissions{},
						Jobs: map[string]*workflow.Job{
							"foo": {},
							"bar": {},
						},
					},
				},
			},
			job:   &workflow.Job{},
			isErr: true,
		},
	}
	p := &policy.JobPermissionsPolicy{}
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
