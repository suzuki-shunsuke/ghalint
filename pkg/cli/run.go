package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func findConfig() string {
	for _, filePath := range []string{"ghalint.yaml", ".ghalint.yaml", "ghalint.yml", ".ghalint.yml"} {
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}
	return ""
}

func readConfig(cfg *Config, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open a configuration file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("parse configuration file as YAML: %w", err)
	}
	return nil
}

func validateConfig(cfg *Config) error {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == "" {
			return errors.New(`policy_name is required`)
		}
		if exclude.PolicyName != "job_secrets" {
			return errors.New(`only the policy "job_secrets" can be excluded`)
		}
		if exclude.WorkflowFilePath == "" {
			return errors.New(`workflow_file_path is required`)
		}
		if exclude.JobName == "" {
			return errors.New(`jobName is required`)
		}
	}
	return nil
}

func (runner *Runner) Run(ctx *cli.Context) error {
	logE := log.New(runner.flags.Version)
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
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		wf := &Workflow{
			FilePath: filePath,
		}
		if err := readWorkflow(filePath, wf); err != nil {
			failed = true
			logerr.WithError(logE, err).Error("read a workflow file")
			continue
		}

		for _, policy := range policies {
			logE := logE.WithField("policy_name", policy.Name())
			if err := policy.Apply(ctx.Context, logE, cfg, wf); err != nil {
				failed = true
				continue
			}
		}
	}
	if failed {
		return errors.New("some workflow files are invalid")
	}
	return nil
}

type Policy interface {
	Name() string
	Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error
}

type Workflow struct {
	FilePath    string `yaml:"-"`
	Jobs        map[string]*Job
	Env         map[string]string
	Permissions map[string]string
}

type Job struct {
	Permissions map[string]string
	Env         map[string]string
	Steps       []interface{}
}
