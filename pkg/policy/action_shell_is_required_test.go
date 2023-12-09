package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestActionShellIsRequiredPolicy_ApplyStep(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		step  *workflow.Step
		isErr bool
	}{
		{
			name: "pass",
			step: &workflow.Step{
				Run:   "echo hello",
				Shell: "bash",
			},
		},
		{
			name:  "step error",
			isErr: true,
			step: &workflow.Step{
				Run: "echo hello",
			},
		},
	}
	p := &policy.ActionShellIsRequiredPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyStep(logE, nil, nil, d.step); err != nil {
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
