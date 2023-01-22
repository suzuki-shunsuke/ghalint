package cli

import (
	"os"

	"gopkg.in/yaml.v3"
)

func readWorkflow(p string, wf *Workflow) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(wf); err != nil {
		return err
	}
	return nil
}
