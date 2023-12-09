package policy

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestJobPermissionsPolicy_ApplyJob(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		jobCtx *JobContext
		job    *workflow.Job
		isErr  bool
	}{
		{
			name: "workflow permissions is empty",
			job:  &workflow.Job{},
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
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
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
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
			jobCtx: &JobContext{
				Workflow: &WorkflowContext{
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
	policy := &JobPermissionsPolicy{}
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
