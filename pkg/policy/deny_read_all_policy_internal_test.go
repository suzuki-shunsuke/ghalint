package policy //nolint:dupl

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestDenyReadAllPermissionPolicy_ApplyJob(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		jobCtx *JobContext
		job    *workflow.Job
		isErr  bool
	}{
		{
			name: "don't use read-all",
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
					Workflow: &workflow.Workflow{},
				},
			},
			job: &workflow.Job{
				Permissions: workflow.NewPermissions(true, false, nil),
			},
			isErr: true,
		},
		{
			name: "job permissions is null and workflow permissions is read-all",
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
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
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
					Workflow: &workflow.Workflow{},
				},
			},
			job: &workflow.Job{
				Permissions: workflow.NewPermissions(false, false, map[string]string{
					"contents": "read",
				}),
			},
		},
	}
	policy := &DenyReadAllPermissionPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := policy.ApplyJob(logE, nil, d.jobCtx, d.job); err != nil {
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
