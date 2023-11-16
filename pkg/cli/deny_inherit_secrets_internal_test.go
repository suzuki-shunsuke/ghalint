package cli

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDenyInheritSecretsPolicy_Apply(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		cfg   *Config
		wf    *Workflow
		isErr bool
	}{
		{
			name: "error",
			wf: &Workflow{
				Jobs: map[string]*Job{
					"release": {
						Secrets: &JobSecrets{
							inherit: true,
						},
					},
				},
			},
			isErr: true,
		},
		{
			name: "pass",
			wf: &Workflow{
				Jobs: map[string]*Job{
					"release": {
						Secrets: &JobSecrets{
							m: map[string]string{
								"foo": "${{secrets.API_KEY}}",
							},
						},
					},
				},
			},
		},
	}
	p := &DenyInheritSecretsPolicy{}
	ctx := context.Background()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
