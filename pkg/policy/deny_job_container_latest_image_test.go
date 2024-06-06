package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestDenyJobContainerLatestImagePolicy_ApplyJob(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		job   *workflow.Job
		isErr bool
	}{
		{
			name: "pass",
			job: &workflow.Job{
				Container: &workflow.Container{
					Image: "node:18",
				},
			},
		},
		{
			name: "job container should have image",
			job: &workflow.Job{
				Container: &workflow.Container{},
			},
			isErr: true,
		},
		{
			name: "job container image should have tag",
			job: &workflow.Job{
				Container: &workflow.Container{
					Image: "node",
				},
			},
			isErr: true,
		},
		{
			name: "latest",
			job: &workflow.Job{
				Container: &workflow.Container{
					Image: "node:latest",
				},
			},
			isErr: true,
		},
	}
	p := &policy.DenyJobContainerLatestImagePolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logE, nil, nil, d.job); err != nil {
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
