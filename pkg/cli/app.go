package cli

import (
	"github.com/urfave/cli/v2"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (flags *LDFlags) AppVersion() string {
	return flags.Version + " (" + flags.Commit + ")"
}

type Runner struct {
	flags *LDFlags
}

func New(flags *LDFlags) *cli.App {
	app := cli.NewApp()
	app.Name = "ghalint"
	app.Usage = "GitHub Actions linter"
	app.Version = flags.AppVersion()
	app.Flags = []cli.Flag{
		// &cli.StringFlag{Name: "log-level", Usage: "log level"},
	}
	runner := &Runner{
		flags: flags,
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
