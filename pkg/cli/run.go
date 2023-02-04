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

func (runner *Runner) Run(ctx *cli.Context) error {
	logE := log.New(runner.flags.Version)

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
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		if runner.validateWorkflow(ctx, logE, cfg, policies, filePath) {
			failed = true
		}
	}
	if failed {
		return errors.New("some workflow files are invalid")
	}
	return nil
}

func (runner *Runner) validateWorkflow(ctx *cli.Context, logE *logrus.Entry, cfg *Config, policies []Policy, filePath string) bool {
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

type Workflow struct {
	FilePath    string `yaml:"-"`
	Jobs        map[string]*Job
	Env         map[string]string
	Permissions *Permissions
}

type Job struct {
	Permissions *Permissions
	Env         map[string]string
	Steps       []interface{}
}

type Permissions struct {
	m        map[string]string
	readAll  bool
	writeAll bool
}

func (permissions *Permissions) Permissions() map[string]string {
	if permissions == nil {
		return nil
	}
	return permissions.m
}

func (permissions *Permissions) ReadAll() bool {
	if permissions == nil {
		return false
	}
	return permissions.readAll
}

func (permissions *Permissions) WriteAll() bool {
	if permissions == nil {
		return false
	}
	return permissions.writeAll
}

func (permissions *Permissions) IsNil() bool {
	if permissions == nil {
		return true
	}
	return permissions.m == nil && !permissions.readAll && !permissions.writeAll
}

func (permissions *Permissions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val interface{}
	if err := unmarshal(&val); err != nil {
		return err
	}
	return convPermissions(val, permissions)
}

func convPermissions(src interface{}, dest *Permissions) error { //nolint:cyclop
	switch p := src.(type) {
	case string:
		switch p {
		case "read-all":
			dest.readAll = true
			return nil
		case "write-all":
			dest.writeAll = true
			return nil
		default:
			return logerr.WithFields(errors.New("unknown permissions"), logrus.Fields{ //nolint:wrapcheck
				"permission": p,
			})
		}
	case map[interface{}]interface{}:
		m := make(map[string]string, len(p))
		for k, v := range p {
			ks, ok := k.(string)
			if !ok {
				return errors.New("permissions key must be string")
			}
			vs, ok := v.(string)
			if !ok {
				return errors.New("permissions value must be string")
			}
			m[ks] = vs
		}
		dest.m = m
		return nil
	case map[string]interface{}:
		m := make(map[string]string, len(p))
		for k, v := range p {
			vs, ok := v.(string)
			if !ok {
				return errors.New("permissions value must be string")
			}
			m[k] = vs
		}
		dest.m = m
		return nil
	default:
		return errors.New("permissions must be map[string]string or 'read-all' or 'write-all'")
	}
}
