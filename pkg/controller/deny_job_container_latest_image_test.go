package controller_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
)

func TestDenyJobContainerLatestImagePolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *controller.Config
		wf    *controller.Workflow
		isErr bool
	}{
		{
			name: "pass",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				Jobs: map[string]*controller.Job{
					"foo": {},
					"bar": {
						Container: &controller.Container{
							Image: "node:18",
						},
					},
				},
			},
		},
		{
			name: "job container should have image",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				Jobs: map[string]*controller.Job{
					"bar": {
						Container: &controller.Container{},
					},
				},
			},
			isErr: true,
		},
		{
			name: "job container image should have tag",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				Jobs: map[string]*controller.Job{
					"bar": {
						Container: &controller.Container{
							Image: "node",
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "latest",
			cfg:  &controller.Config{},
			wf: &controller.Workflow{
				Jobs: map[string]*controller.Job{
					"bar": {
						Container: &controller.Container{
							Image: "node:latest",
						},
					},
				},
			},
			isErr: true,
		},
	}
	policy := &controller.DenyJobContainerLatestImagePolicy{}
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := policy.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
