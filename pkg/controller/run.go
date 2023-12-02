package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) Run(ctx context.Context, logE *logrus.Entry, cfgFilePath string) error {
	cfg := &config.Config{}
	if err := c.readConfig(cfg, cfgFilePath); err != nil {
		return err
	}

	filePaths, err := workflow.List(c.fs)
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	policies := []Policy{
		&policy.JobPermissionsPolicy{},
		policy.NewWorkflowSecretsPolicy(),
		policy.NewJobSecretsPolicy(),
		&policy.DenyReadAllPermissionPolicy{},
		&policy.DenyWriteAllPermissionPolicy{},
		&policy.DenyInheritSecretsPolicy{},
		&policy.DenyJobContainerLatestImagePolicy{},
		&policy.GitHubAppShouldLimitRepositoriesPolicy{},
		&policy.GitHubAppShouldLimitPermissionsPolicy{},
		policy.NewActionRefShouldBeSHA1Policy(),
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		if c.validateWorkflow(ctx, logE, cfg, policies, filePath) {
			failed = true
		}
	}
	if failed {
		return debugError(errors.New("some workflow files are invalid"))
	}
	return nil
}

func (c *Controller) validateWorkflow(ctx context.Context, logE *logrus.Entry, cfg *config.Config, policies []Policy, filePath string) bool {
	wf := &workflow.Workflow{
		FilePath: filePath,
	}
	if err := workflow.Read(c.fs, filePath, wf); err != nil {
		logerr.WithError(logE, err).Error("read a workflow file")
		return true
	}

	failed := false
	for _, policy := range policies {
		logE := logE.WithFields(logrus.Fields{
			"policy_name": policy.Name(),
			"reference":   fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", policy.ID()),
		})
		if err := policy.Apply(ctx, logE, cfg, wf); err != nil {
			failed = true
			continue
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
