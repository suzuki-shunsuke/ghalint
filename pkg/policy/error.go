package policy

import "errors"

var (
	errPermissionHyphenIsRequired = errors.New("an input `permission-*` is required")
	errPermissionsIsRequired      = errors.New("the input `permissions` is required")
	errRepositoriesIsRequired     = errors.New("the input `repositories` is required")
	errEmpty                      = errors.New("")
)
