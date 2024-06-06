package policy_test //nolint:dupl

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestDenyWriteAllPermissionPolicy_ApplyJob(t *testing.T) {
	t.Parallel()
	data := []struct {
		name   string
		jobCtx *policy.JobContext
		job    *workflow.Job
		isErr  bool
	}{
		{
			name: "don't use write-all",
			job: &workflow.Job{
				Permissions: workflow.NewPermissions(false, true, nil),
			},
			isErr: true,
		},
		{
			name: "job permissions is null and workflow permissions is write-all",
			jobCtx: &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					Workflow: &workflow.Workflow{
						Permissions: workflow.NewPermissions(false, true, nil),
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
					"contents": "write",
				}),
			},
		},
	}
	p := &policy.DenyWriteAllPermissionPolicy{}
	logE := logrus.NewEntry(logrus.New())
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
			if err := p.ApplyJob(logE, nil, d.jobCtx, d.job); err != nil {
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
