package controller

import (
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

type WorkflowPolicy interface {
	Name() string
	ID() string
	ApplyWorkflow(logger *slog.Logger, cfg *config.Config, wfCtx *policy.WorkflowContext, wf *workflow.Workflow) error
}

type JobPolicy interface {
	Name() string
	ID() string
	ApplyJob(logger *slog.Logger, cfg *config.Config, jobCtx *policy.JobContext, job *workflow.Job) error
}

type StepPolicy interface {
	Name() string
	ID() string
	ApplyStep(logger *slog.Logger, cfg *config.Config, stepCtx *policy.StepContext, step *workflow.Step) error
}
