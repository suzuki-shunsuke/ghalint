package policy //nolint:dupl

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestDenyReadAllPermissionPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    string
		isErr bool
	}{
		{
			name: "don't use read-all",
			cfg:  &config.Config{},
			wf: `jobs:
  foo:
    permissions: read-all`,
			isErr: true,
		},
		{
			name: "job permissions is null and workflow permissions is read-all",
			cfg:  &config.Config{},
			wf: `permissions: read-all
jobs:
  foo: {}`,
			isErr: true,
		},
		{
			name: "pass",
			cfg:  &config.Config{},
			wf: `jobs:
  foo:
    permissions:
      contents: read`,
		},
	}
	policy := &DenyReadAllPermissionPolicy{}
	logE := logrus.NewEntry(logrus.New())
	ctx := context.Background()
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			wf := &workflow.Workflow{}
			if err := yaml.Unmarshal([]byte(d.wf), wf); err != nil {
				t.Fatal(err)
			}
			if err := policy.Apply(ctx, logE, d.cfg, wf); err != nil {
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
