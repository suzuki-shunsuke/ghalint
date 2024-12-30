package policy

import (
	"errors"
	"fmt"
)

var (
	errPermissionsIsRequired  = errors.New("the input `permissions` is required")
	errRepositoriesIsRequired = errors.New("the input `repositories` is required")
	errEmpty                  = errors.New("")
)

type Info struct {
	Name    string `json:"name,omitempty"`
	ID      string `json:"id,omitempty"`
	URL     string `json:"url,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorInfo struct {
	Policy        *Info             `json:"policy,omitempty"`
	FilePath      string            `json:"file_path,omitempty"`
	JobName       string            `json:"job_name,omitempty"`
	StepName      string            `json:"step_name,omitempty"`
	StepID        string            `json:"step_id,omitempty"`
	StepIndex     int               `json:"step_index,omitempty"`
	ActionName    string            `json:"action_name,omitempty"`
	ActionVersion string            `json:"action_version,omitempty"`
	Message       string            `json:"message,omitempty"`
	With          map[string]string `json:"with,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
}

func GetURL(id string) string {
	return fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", id)
}
