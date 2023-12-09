package act

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) Run(_ context.Context, logE *logrus.Entry, cfgFilePath string, args ...string) error {
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
		policy.NewActionRefShouldBeSHA1Policy(),
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("action_file_path", filePath)
		if c.validateAction(logE, cfg, stepPolicies, filePath) {
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
	for _, file := range []string{"action.yaml", "action.yml"} {
		f, err := afero.Exists(c.fs, file)
		if err != nil {
			return nil, fmt.Errorf("check if the action file exists: %w", err)
		}
		if f {
			return []string{file}, nil
		}
	}
	return nil, nil
}

func (c *Controller) validateAction(logE *logrus.Entry, cfg *config.Config, stepPolicies []controller.StepPolicy, filePath string) bool {
	action := &workflow.Action{}
	if err := workflow.ReadAction(c.fs, filePath, action); err != nil {
		logerr.WithError(logE, err).Error("read an action file")
		return true
	}

	stepCtx := &policy.StepContext{
		FilePath: filePath,
		Action:   action,
	}

	failed := false

	if c.applyStepPolicies(logE, cfg, stepCtx, action, stepPolicies) {
		failed = true
	}

	return failed
}

type Policy interface {
	Name() string
	ID() string
}

func withPolicyReference(logE *logrus.Entry, p Policy) *logrus.Entry {
	return logE.WithFields(logrus.Fields{
		"policy_name": p.Name(),
		"reference":   fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", p.ID()),
	})
}

func (c *Controller) applyStepPolicies(logE *logrus.Entry, cfg *config.Config, stepCtx *policy.StepContext, action *workflow.Action, stepPolicies []controller.StepPolicy) bool {
	failed := false
	for _, stepPolicy := range stepPolicies {
		logE := withPolicyReference(logE, stepPolicy)
		if c.applyStepPolicy(logE, cfg, stepCtx, action, stepPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyStepPolicy(logE *logrus.Entry, cfg *config.Config, stepCtx *policy.StepContext, action *workflow.Action, stepPolicy controller.StepPolicy) bool {
	failed := false
	for _, step := range action.Runs.Steps {
		logE := logE
		if step.ID != "" {
			logE = logE.WithField("step_id", step.ID)
		}
		if step.Name != "" {
			logE = logE.WithField("step_name", step.Name)
		}
		if err := stepPolicy.ApplyStep(logE, cfg, stepCtx, step); err != nil {
			if err.Error() != "" {
				logerr.WithError(logE, err).Error("the step violates policies")
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
			return fmt.Errorf("read a configuration file: %w", logerr.WithFields(err, logrus.Fields{
				"config_file": cfgFilePath,
			}))
		}
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validate a configuration file: %w", logerr.WithFields(err, logrus.Fields{
				"config_file": cfgFilePath,
			}))
		}
	}
	return nil
}
