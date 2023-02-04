package cli //nolint:dupl

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDenyReadAllPermissionPolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *Config
		wf    *Workflow
		isErr bool
	}{
		{
			name: "don't use read-all",
			cfg:  &Config{},
			wf: &Workflow{
				Permissions: &Permissions{
					m: map[string]string{},
				},
				Jobs: map[string]*Job{
					"foo": {
						Permissions: &Permissions{
							readAll: true,
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "job permissions is null and workflow permissions is read-all",
			cfg:  &Config{},
			wf: &Workflow{
				Permissions: &Permissions{
					readAll: true,
				},
				Jobs: map[string]*Job{
					"foo": {},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &Config{},
			wf: &Workflow{
				Jobs: map[string]*Job{
					"foo": {
						Permissions: &Permissions{
							m: map[string]string{
								"contents": "read",
							},
						},
					},
				},
			},
		},
	}
	policy := &DenyReadAllPermissionPolicy{}
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
