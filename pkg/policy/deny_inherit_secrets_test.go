package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestDenyInheritSecretsPolicy_ApplyJob(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		job   string
		isErr bool
	}{
		{
			name:  "error",
			job:   `secrets: inherit`,
			isErr: true,
		},
		{
			name: "pass",
			job: `secrets:
      foo: ${{secrets.API_KEY}}`,
		},
	}
	p := &policy.DenyInheritSecretsPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			job := &workflow.Job{}
			if err := yaml.Unmarshal([]byte(d.job), job); err != nil {
				t.Fatal(err)
			}
			if err := p.ApplyJob(logE, nil, nil, job); err != nil {
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
