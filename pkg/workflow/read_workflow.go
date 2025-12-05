package workflow

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"gopkg.in/yaml.v3"
)

func Read(fs afero.Fs, p string, wf *Workflow) error {
	f, err := fs.Open(p)
	if err != nil {
		return fmt.Errorf("open a workflow file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(wf); err != nil {
		err := fmt.Errorf("parse a workflow file as YAML: %w", err)
		if errors.Is(err, io.EOF) {
			return slogerr.With(err, "reference", "https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/codes/001.md") //nolint:wrapcheck
		}
		return err
	}
	return nil
}
