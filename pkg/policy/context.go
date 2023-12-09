package policy

import "github.com/suzuki-shunsuke/ghalint/pkg/workflow"

type WorkflowContext struct {
	FilePath string
	Workflow *workflow.Workflow
}

type JobContext struct {
	Workflow *WorkflowContext
	Name     string
}
