package policy

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestJobPermissionsPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    string
		isErr bool
	}{
		{
			name: "workflow permissions is empty",
			cfg:  &config.Config{},
			wf: `permissions: {}
jobs:
  foo: {}
  bar: {}`,
		},
		{
			name: "workflow has only one job",
			cfg:  &config.Config{},
			wf: `permissions:
  contents: read
jobs:
  foo: {}`,
		},
		{
			name: "job should have permissions",
			cfg:  &config.Config{},
			wf: `jobs:
  foo: {}
  bar: {}`,
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
