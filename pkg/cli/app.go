package cli

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/urfave-cli-v3-help-all/helpall"
	"github.com/urfave/cli/v3"
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
	logE  *logrus.Entry
}

func New(flags *LDFlags, fs afero.Fs, logE *logrus.Entry) *cli.Command {
	app := &cli.Command{}
	app.Name = "ghalint"
	app.Usage = "GitHub Actions linter"
	app.Version = flags.AppVersion()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log-color",
			Usage: "log color. auto(default)|always|never",
			Sources: cli.EnvVars(
				"GHALINT_LOG_COLOR",
			),
		},
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "log level",
			Sources: cli.EnvVars(
				"GHALINT_LOG_LEVEL",
			),
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "configuration file path",
			Sources: cli.EnvVars(
				"GHALINT_CONFIG",
			),
		},
	}
	runner := &Runner{
		flags: flags,
		fs:    fs,
		logE:  logE,
	}
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "lint GitHub Actions Workflows",
			Action: runner.Run,
			Flags:  []cli.Flag{},
		},
		{
			Name: "run-action",
			Aliases: []string{
				"act",
			},
			Usage:  "lint actions",
			Action: runner.RunAction,
			Flags:  []cli.Flag{},
		},
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(_ context.Context, cmd *cli.Command) error {
				cli.ShowVersion(cmd)
				return nil
			},
		},
	}
	return helpall.With(app, nil)
}
