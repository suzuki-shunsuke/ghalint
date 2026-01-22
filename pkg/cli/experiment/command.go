package experiment

import (
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment/validateinput"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

func New(logger *slogutil.Logger, fs afero.Fs, validateInputArgs *validateinput.Args) *cli.Command {
	return &cli.Command{
		Name:        "experiment",
		Aliases:     []string{"exp"},
		Usage:       "experimental commands",
		Description: "experimental commands. These commands are not stable and may change in the future without major updates.",
		Commands: []*cli.Command{
			validateinput.New(logger, fs, validateInputArgs),
		},
	}
}
