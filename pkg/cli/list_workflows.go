package cli

import "path/filepath"

func listWorkflows() ([]string, error) {
	files, err := filepath.Glob(".github/workflows/*.yml")
	if err != nil {
		return nil, err
	}
	files2, err := filepath.Glob(".github/workflows/*.yaml")
	if err != nil {
		return nil, err
	}
	return append(files, files2...), nil
}
