package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type Controller struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Controller {
	return &Controller{
		fs: fs,
	}
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
