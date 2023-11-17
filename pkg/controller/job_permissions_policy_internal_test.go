package controller

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestJobPermissionsPolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *Config
		wf    *Workflow
		isErr bool
	}{
		{
			name: "workflow permissions is empty",
			cfg:  &Config{},
			wf: &Workflow{
				Permissions: &Permissions{
					m: map[string]string{},
				},
				Jobs: map[string]*Job{
					"foo": {},
					"bar": {},
				},
			},
		},
		{
			name: "workflow has only one job",
			cfg:  &Config{},
			wf: &Workflow{
				Permissions: &Permissions{
					m: map[string]string{
						"contents": "read",
					},
				},
				Jobs: map[string]*Job{
					"foo": {},
				},
			},
		},
		{
			name: "job should have permissions",
			cfg:  &Config{},
			wf: &Workflow{
				Jobs: map[string]*Job{
					"foo": {},
					"bar": {},
				},
			},
			isErr: true,
		},
	}
	policy := &JobPermissionsPolicy{}
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
