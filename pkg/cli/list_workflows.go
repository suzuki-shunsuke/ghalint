package cli

import (
	"fmt"

	"github.com/spf13/afero"
)

func listWorkflows(fs afero.Fs) ([]string, error) {
	files, err := afero.Glob(fs, ".github/workflows/*.yml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yml: %w", err)
	}
	files2, err := afero.Glob(fs, ".github/workflows/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("find .github/workflows/*.yaml: %w", err)
	}
	return append(files, files2...), nil
}
