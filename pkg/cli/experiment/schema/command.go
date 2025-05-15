package schema

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment/schema/validate"
	"github.com/urfave/cli/v3"
)

func New(logE *logrus.Entry, fs afero.Fs) *cli.Command {
	return &cli.Command{
		Name:        "schema",
		Usage:       "schema validation",
		Description: "schema validation",
		Commands: []*cli.Command{
			validate.New(logE, fs),
		},
	}
}
