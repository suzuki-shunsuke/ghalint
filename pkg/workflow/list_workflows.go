package workflow

import (
	"fmt"
	"path/filepath"
)

func List() ([]string, error) {
	files, err := filepath.Glob(".github/workflows/*.yml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yml: %w", err)
	}
	files2, err := filepath.Glob(".github/workflows/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yaml: %w", err)
	}
	return append(files, files2...), nil
}
