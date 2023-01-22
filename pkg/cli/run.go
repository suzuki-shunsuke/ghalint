package cli

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) Run(ctx *cli.Context) error {
	filePaths, err := listWorkflows()
	if err != nil {
		return err
	}
	logE := log.New(runner.flags.Version)
	policies := []Policy{
		&JobPermissionsPolicy{},
		NewWorkflowSecretsPolicy(),
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
			if err := policy.Apply(ctx.Context, logE, wf); err != nil {
				failed = true
				logerr.WithError(logE, err).Error("apply a policy")
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
	Apply(ctx context.Context, logE *logrus.Entry, wf *Workflow) error
}

type Workflow struct {
	FilePath string `yaml:"-"`
	Jobs     map[string]*Job
	Env      map[string]string
}

type Job struct {
	Permissions map[string]string
}
