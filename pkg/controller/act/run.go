package act

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/action"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Controller) Run(_ context.Context, logger *slog.Logger, cfgFilePath string, args ...string) error {
	cfg := &config.Config{}
	if err := c.readConfig(cfg, cfgFilePath); err != nil {
		return err
	}

	filePaths, err := c.listFiles(args...)
	if err != nil {
		return fmt.Errorf("find action files: %w", err)
	}
	stepPolicies := []controller.StepPolicy{
		&policy.GitHubAppShouldLimitRepositoriesPolicy{},
		&policy.GitHubAppShouldLimitPermissionsPolicy{},
		&policy.ActionShellIsRequiredPolicy{},
		policy.NewActionRefShouldBeSHAPolicy(),
		&policy.CheckoutPersistCredentialShouldBeFalsePolicy{},
	}
	failed := false
	for _, filePath := range filePaths {
		logger := logger.With("action_file_path", filePath)
		if c.validateAction(logger, cfg, stepPolicies, filePath) {
			failed = true
		}
	}
	if failed {
		return debugError(errors.New("some action files are invalid"))
	}
	return nil
}

func (c *Controller) listFiles(args ...string) ([]string, error) {
	if len(args) != 0 {
		return args, nil
	}

	return action.Find(c.fs) //nolint:wrapcheck
}

func (c *Controller) validateAction(logger *slog.Logger, cfg *config.Config, stepPolicies []controller.StepPolicy, filePath string) bool {
	action := &workflow.Action{}
	if err := workflow.ReadAction(c.fs, filePath, action); err != nil {
		slogerr.WithError(logger, err).Error("read an action file")
		return true
	}

	stepCtx := &policy.StepContext{
		FilePath: filePath,
		Action:   action,
	}

	return c.applyStepPolicies(logger, cfg, stepCtx, action, stepPolicies)
}

type Policy interface {
	Name() string
	ID() string
}

func withPolicyReference(logger *slog.Logger, p Policy) *slog.Logger {
	return logger.With(
		"policy_name", p.Name(),
		"reference", fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", p.ID()),
	)
}

func (c *Controller) applyStepPolicies(logger *slog.Logger, cfg *config.Config, stepCtx *policy.StepContext, action *workflow.Action, stepPolicies []controller.StepPolicy) bool {
	failed := false
	for _, stepPolicy := range stepPolicies {
		logger := withPolicyReference(logger, stepPolicy)
		if c.applyStepPolicy(logger, cfg, stepCtx, action, stepPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyStepPolicy(logger *slog.Logger, cfg *config.Config, stepCtx *policy.StepContext, action *workflow.Action, stepPolicy controller.StepPolicy) bool {
	failed := false
	for _, step := range action.Runs.Steps {
		logger := logger
		if step.ID != "" {
			logger = logger.With("step_id", step.ID)
		}
		if step.Name != "" {
			logger = logger.With("step_name", step.Name)
		}
		if err := stepPolicy.ApplyStep(logger, cfg, stepCtx, step); err != nil {
			if err.Error() != "" {
				slogerr.WithError(logger, err).Error("the step violates policies")
			}
			failed = true
		}
	}
	return failed
}

func (c *Controller) readConfig(cfg *config.Config, cfgFilePath string) error {
	if cfgFilePath == "" {
		if c := config.Find(c.fs); c != "" {
			cfgFilePath = c
		}
	}
	if cfgFilePath != "" {
		if err := config.Read(c.fs, cfg, cfgFilePath); err != nil {
			return fmt.Errorf("read a configuration file: %w", slogerr.With(err,
				"config_file", cfgFilePath,
			))
		}
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validate a configuration file: %w", slogerr.With(err,
				"config_file", cfgFilePath,
			))
		}
		config.ConvertPath(cfg)
	}
	return nil
}
