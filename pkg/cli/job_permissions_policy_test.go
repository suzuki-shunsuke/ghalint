package cli_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli"
)

func TestJobPermissionsPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name string
		cfg  *cli.Config
		wf   *cli.Workflow
		exp  bool
	}{
		{
			name: "workflow permissions is empty",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Permissions: map[string]string{},
				Jobs: map[string]*cli.Job{
					"foo": {},
					"bar": {},
				},
			},
		},
		{
			name: "workflow has only one job",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Permissions: map[string]string{
					"contents": "read",
				},
				Jobs: map[string]*cli.Job{
					"foo": {},
				},
			},
		},
		{
			name: "job should have permissions",
			cfg:  &cli.Config{},
			wf: &cli.Workflow{
				Jobs: map[string]*cli.Job{
					"foo": {},
					"bar": {},
				},
			},
			exp: false,
		},
	}
	policy := &cli.JobPermissionsPolicy{}
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := policy.Apply(ctx, logE, d.cfg, d.wf); err != nil {
				if d.exp {
					t.Fatal(err)
				}
				return
			}
			if d.exp {
				t.Fatal("error must be returned")
			}
		})
	}
}
