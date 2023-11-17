package controller

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Controller struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Controller {
	return &Controller{
		fs: fs,
	}
}

func (c *Controller) Run(ctx context.Context, logE *logrus.Entry) error {
	cfg := &Config{}
	if cfgFilePath := findConfig(c.fs); cfgFilePath != "" {
		if err := readConfig(c.fs, cfg, cfgFilePath); err != nil {
			logE.WithError(err).Error("read a configuration file")
			return err
		}
	}
	if err := validateConfig(cfg); err != nil {
		logE.WithError(err).Error("validate a configuration file")
		return err
	}
	filePaths, err := listWorkflows(c.fs)
	if err != nil {
		logE.Error(err)
		return err
	}
	policies := []Policy{
		&JobPermissionsPolicy{},
		NewWorkflowSecretsPolicy(),
		NewJobSecretsPolicy(),
		&DenyReadAllPermissionPolicy{},
		&DenyWriteAllPermissionPolicy{},
		&DenyInheritSecretsPolicy{},
		&DenyJobContainerLatestImagePolicy{},
		NewActionRefShouldBeSHA1Policy(),
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		if c.validateWorkflow(ctx, logE, cfg, policies, filePath) {
			failed = true
		}
	}
	if failed {
		return errors.New("some workflow files are invalid")
	}
	return nil
}

func (c *Controller) validateWorkflow(ctx context.Context, logE *logrus.Entry, cfg *Config, policies []Policy, filePath string) bool {
	wf := &Workflow{
		FilePath: filePath,
	}
	if err := readWorkflow(c.fs, filePath, wf); err != nil {
		logerr.WithError(logE, err).Error("read a workflow file")
		return true
	}

	failed := false
	for _, policy := range policies {
		logE := logE.WithField("policy_name", policy.Name())
		if err := policy.Apply(ctx, logE, cfg, wf); err != nil {
			failed = true
			continue
		}
	}
	return failed
}

type Policy interface {
	Name() string
	Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error
}
