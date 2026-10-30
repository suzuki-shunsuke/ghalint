package action

import (
	"fmt"
	"path/filepath"
)

func Find() ([]string, error) {
	patterns := []string{
		"action.yaml",
		"action.yml",
		"*/action.yaml",
		"*/action.yml",
		"*/*/action.yaml",
		"*/*/action.yml",
		"*/*/*/action.yaml",
		"*/*/*/action.yml",
	}

	files := []string{}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("check if the action file exists: %w", err)
		}
		files = append(files, matches...)
	}
	return files, nil
}
