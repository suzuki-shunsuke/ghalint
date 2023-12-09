package policy

import "errors"

var (
	errPermissionsIsRequired  = errors.New("the input `permissions` is required")
	errRepositoriesIsRequired = errors.New("the input `repositories` is required")
	errEmpty                  = errors.New("")
)
