package cli

import (
	"context"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
)

func (r *Runner) Run(ctx *cli.Context) error {
	logE := log.New(r.flags.Version)

	if color := os.Getenv("GHALINT_LOG_COLOR"); color != "" {
		log.SetColor(color, logE)
	}

	cfg := &Config{}
	if cfgFilePath := findConfig(); cfgFilePath != "" {
		if err := readConfig(cfg, cfgFilePath); err != nil {
			logE.WithError(err).Error("read a configuration file")
			return err
		}
	}
	if err := validateConfig(cfg); err != nil {
		logE.WithError(err).Error("validate a configuration file")
		return err
	}
	filePaths, err := listWorkflows()
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
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		if r.validateWorkflow(ctx, logE, cfg, policies, filePath) {
			failed = true
		}
	}
	if failed {
		return errors.New("some workflow files are invalid")
	}
	return nil
}

func (r *Runner) validateWorkflow(ctx *cli.Context, logE *logrus.Entry, cfg *Config, policies []Policy, filePath string) bool {
	wf := &Workflow{
		FilePath: filePath,
	}
	if err := readWorkflow(filePath, wf); err != nil {
		logerr.WithError(logE, err).Error("read a workflow file")
		return true
	}

	failed := false
	for _, policy := range policies {
		logE := logE.WithField("policy_name", policy.Name())
		if err := policy.Apply(ctx.Context, logE, cfg, wf); err != nil {
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
