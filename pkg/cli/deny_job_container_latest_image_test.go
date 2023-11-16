package cli_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
)

func TestDenyJobContainerLatestImagePolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *cli.Config
		wf    *cli.Workflow
		isErr bool
	}{
		{
			name: "pass",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"foo": {},
					"bar": {
						Container: &cli.Container{
							Image: "node:18",
						},
					},
				},
			},
		},
		{
			name: "job container should have image",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"bar": {
						Container: &cli.Container{},
					},
				},
			},
			isErr: true,
		},
		{
			name: "job container image should have tag",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"bar": {
						Container: &cli.Container{
							Image: "node",
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "latest",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"bar": {
						Container: &cli.Container{
							Image: "node:latest",
						},
					},
				},
			},
			isErr: true,
		},
	}
	policy := &cli.DenyJobContainerLatestImagePolicy{}
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
