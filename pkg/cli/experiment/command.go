package experiment

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment/schema"
	"github.com/urfave/cli/v3"
)

func New(logE *logrus.Entry, fs afero.Fs) *cli.Command {
	return &cli.Command{
		Name:        "experiment",
		Aliases:     []string{"exp"},
		Usage:       "experimental commands",
		Description: "experimental commands. These commands are not stable and may change in the future without major updates.",
		Commands: []*cli.Command{
			schema.New(logE, fs),
		},
	}
}
