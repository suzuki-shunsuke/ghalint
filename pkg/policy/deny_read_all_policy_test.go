package policy_test //nolint:dupl

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestDenyReadAllPermissionPolicy_ApplyJob(t *testing.T) {
	t.Parallel()
	data := []struct {
		name   string
		jobCtx *policy.JobContext
		job    *workflow.Job
		isErr  bool
	}{
		{
			name: "don't use read-all",
			job: &workflow.Job{
				Permissions: workflow.NewPermissions(true, false, nil),
			},
			isErr: true,
		},
		{
			name: "job permissions is null and workflow permissions is read-all",
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{
						Permissions: workflow.NewPermissions(true, false, nil),
					},
				},
			},
			job:   &workflow.Job{},
			isErr: true,
		},
		{
			name: "pass",
			job: &workflow.Job{
				Permissions: workflow.NewPermissions(false, false, map[string]string{
					"contents": "read",
				}),
			},
		},
	}
	p := &policy.DenyReadAllPermissionPolicy{}
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		if d.jobCtx == nil {
			d.jobCtx = &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{},
				},
			}
		}
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logger, nil, d.jobCtx, d.job); err != nil {
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
