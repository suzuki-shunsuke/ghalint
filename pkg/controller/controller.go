package controller

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type Controller struct {
	fs     afero.Fs
	stderr io.Writer
}

func New(fs afero.Fs, input *InputNew) *Controller {
	return &Controller{
		fs:     fs,
		stderr: input.Stderr,
	}
}

type InputNew struct {
	Stderr io.Writer
}

type WorkflowPolicy interface {
	Name() string
	ID() string
	ApplyWorkflow(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, wf *workflow.Workflow) error
}

type JobPolicy interface {
	Name() string
	ID() string
	ApplyJob(logE *logrus.Entry, cfg *config.Config, jobCtx *policy.JobContext, job *workflow.Job) error
}

type StepPolicy interface {
	Name() string
	ID() string
	ApplyStep(logE *logrus.Entry, cfg *config.Config, stepCtx *policy.StepContext, step *workflow.Step) error
}
