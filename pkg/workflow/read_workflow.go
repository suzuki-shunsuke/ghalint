package workflow

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func Read(fs afero.Fs, p string, wf *Workflow) error {
	f, err := fs.Open(p)
	if err != nil {
		return fmt.Errorf("open a workflow file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(wf); err != nil {
		return fmt.Errorf("parse a workflow file as YAML: %w", err)
	}
	return nil
}
