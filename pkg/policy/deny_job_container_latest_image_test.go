package policy_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestDenyJobContainerLatestImagePolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    *workflow.Workflow
		isErr bool
	}{
		{
			name: "pass",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"foo": {},
					"bar": {
						Container: &workflow.Container{
							Image: "node:18",
						},
					},
				},
			},
		},
		{
			name: "job container should have image",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"bar": {
						Container: &workflow.Container{},
					},
				},
			},
			isErr: true,
		},
		{
			name: "job container image should have tag",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"bar": {
						Container: &workflow.Container{
							Image: "node",
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "latest",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				Jobs: map[string]*workflow.Job{
					"bar": {
						Container: &workflow.Container{
							Image: "node:latest",
						},
					},
				},
			},
			isErr: true,
		},
	}
	p := &policy.DenyJobContainerLatestImagePolicy{}
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
