package cli

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/log"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) Run(ctx *cli.Context) error {
	// find files .github/workflows/*.ya?ml
	filePaths, err := listWorkflows()
	if err != nil {
		return err
	}
	logE := log.New(runner.flags.Version)
	workflowSecretsPolicy, err := NewWorkflowSecretsPolicy()
	if err != nil {
		return err
	}
	policies := []Policy{
		&JobPermissionsPolicy{},
		workflowSecretsPolicy,
	}
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		wf := &Workflow{
			FilePath: filePath,
		}
		if err := readWorkflow(filePath, wf); err != nil {
			return err
		}
		// apply policies
		for _, policy := range policies {
			logE := logE.WithField("policy_name", policy.Name())
			if err := policy.Apply(ctx.Context, logE, wf); err != nil {
				return err
			}
		}
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
