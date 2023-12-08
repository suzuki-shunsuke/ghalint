package policy

import "errors"

var (
	errPermissionsIsRequired   = errors.New("the input `permissions` is required")
	errRepositoriesIsRequired  = errors.New("the input `repositories` is required")
	errJobViolatePolicy      = errors.New("the job violates policies")
	errWorkflowViolatePolicy = errors.New("the workflow violates policies")
)
