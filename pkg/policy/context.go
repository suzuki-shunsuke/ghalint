package policy

import "github.com/suzuki-shunsuke/ghalint/pkg/workflow"

type WorkflowContext struct {
	FilePath string
	Workflow *workflow.Workflow
	Content  []byte
}

type JobContext struct {
	Name     string
	Workflow *WorkflowContext
	Job      *workflow.Job
}

type StepContext struct {
	FilePath string
	Action   *workflow.Action
	Job      *JobContext
}
