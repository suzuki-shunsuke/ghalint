package policy

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestDenyInheritSecretsPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    string
		isErr bool
	}{
		{
			name: "error",
			wf: `jobs:
  releases:
    secrets: inherit`,
			isErr: true,
		},
		{
			name: "pass",
			wf: `jobs:
  release:
    secrets:
      foo: ${{secrets.API_KEY}}`,
		},
	}
	p := &DenyInheritSecretsPolicy{}
	ctx := context.Background()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			wf := &workflow.Workflow{}
			if err := yaml.Unmarshal([]byte(d.wf), wf); err != nil {
				t.Fatal(err)
			}
			if err := p.Apply(ctx, logE, d.cfg, wf); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
