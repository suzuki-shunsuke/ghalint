package policy

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type ActionShellIsRequiredPolicy struct{}

func (p *ActionShellIsRequiredPolicy) Name() string {
	return "action_shell_is_required"
}

func (p *ActionShellIsRequiredPolicy) ID() string {
	return "011"
}

func (p *ActionShellIsRequiredPolicy) ApplyStep(_ *logrus.Entry, _ *config.Config, _ *StepContext, step *workflow.Step) error {
	if step.Run != "" && step.Shell == "" {
		return errors.New("shell is required if run is set")
	}
	return nil
}
