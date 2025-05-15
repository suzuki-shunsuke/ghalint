package action

import (
	"fmt"

	"github.com/spf13/afero"
)

func Find(fs afero.Fs) ([]string, error) {
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
		matches, err := afero.Glob(fs, pattern)
		if err != nil {
			return nil, fmt.Errorf("check if the action file exists: %w", err)
		}
		files = append(files, matches...)
	}
	return files, nil
}
