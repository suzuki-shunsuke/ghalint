package policy

import (
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type WorkflowContext struct {
	FilePath string
	Workflow *workflow.Workflow
	// Lockfile is nil if the lockfile doesn't exist or ghalint fails to read it.
	Lockfile *lockfile.Lockfile
	// LockfileInvalid is true if ghalint fails to read the lockfile.
	LockfileInvalid bool
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
