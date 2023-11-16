package cli

import (
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (f *LDFlags) AppVersion() string {
	return f.Version + " (" + f.Commit + ")"
}

type Runner struct {
	flags *LDFlags
	fs    afero.Fs
}

func New(flags *LDFlags, fs afero.Fs) *cli.App {
	app := cli.NewApp()
	app.Name = "ghalint"
	app.Usage = "GitHub Actions linter"
	app.Version = flags.AppVersion()
	app.Flags = []cli.Flag{
		// &cli.StringFlag{Name: "log-level", Usage: "log level"},
	}
	runner := &Runner{
		flags: flags,
		fs:    fs,
	}
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "lint GitHub Actions Workflows",
			Action: runner.Run,
			Flags:  []cli.Flag{},
		},
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(ctx *cli.Context) error {
				cli.ShowVersion(ctx)
				return nil
			},
		},
	}
	return app
}
