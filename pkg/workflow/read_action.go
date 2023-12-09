package workflow

import (
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
)

func ReadAction(fs afero.Fs, p string, action *Action) error {
	f, err := fs.Open(p)
	if err != nil {
		return fmt.Errorf("open an action file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(action); err != nil {
		err := fmt.Errorf("parse an action file as YAML: %w", err)
		if errors.Is(err, io.EOF) {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"reference": "https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/codes/001.md",
			})
		}
		return err
	}
	return nil
}
