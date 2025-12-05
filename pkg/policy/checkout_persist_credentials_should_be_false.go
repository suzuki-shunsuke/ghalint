package policy

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type CheckoutPersistCredentialShouldBeFalsePolicy struct{}

func (p *CheckoutPersistCredentialShouldBeFalsePolicy) Name() string {
	return "checkout_persist_credentials_should_be_false"
}

func (p *CheckoutPersistCredentialShouldBeFalsePolicy) ID() string {
	return "013"
}

func (p *CheckoutPersistCredentialShouldBeFalsePolicy) ApplyStep(_ *slog.Logger, cfg *config.Config, stepCtx *StepContext, step *workflow.Step) error {
	if p.excluded(stepCtx, cfg.Excludes) {
		return nil
	}
	if !strings.HasPrefix(step.Uses, "actions/checkout@") {
		return nil
	}
	f, ok := step.With["persist-credentials"]
	if !ok {
		return errors.New("persist-credentials should be false")
	}
	if f != "false" {
		return errors.New("persist-credentials should be false")
	}
	return nil
}

func (p *CheckoutPersistCredentialShouldBeFalsePolicy) excluded(stepCtx *StepContext, excludes []*config.Exclude) bool {
	for _, exclude := range excludes {
		if exclude.PolicyName != p.Name() {
			continue
		}
		if stepCtx.Action != nil {
			if exclude.ActionFilePath != stepCtx.FilePath {
				continue
			}
			return true
		}
		if exclude.JobName != stepCtx.Job.Name {
			continue
		}
		if exclude.WorkflowFilePath != stepCtx.Job.Workflow.FilePath {
			continue
		}
		return true
	}
	return false
}
