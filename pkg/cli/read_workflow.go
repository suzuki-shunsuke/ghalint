package cli

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func readWorkflow(p string, wf *Workflow) error {
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("open a workflow file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(wf); err != nil {
		return fmt.Errorf("parse a workflow file as YAML: %w", err)
	}
	return nil
}
