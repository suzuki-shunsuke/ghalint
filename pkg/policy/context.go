package policy

import (
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type WorkflowContext struct {
	FilePath string
	Workflow *workflow.Workflow
	// Lockfile is nil if the lockfile doesn't exist.
	// If ghalint fails to read the lockfile, action_ref_should_be_full_length_commit_sha isn't applied.
	Lockfile *lockfile.Lockfile
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
