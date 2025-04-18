package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/helpall"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/vcmd"
	"github.com/urfave/cli/v3"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

type Runner struct {
	flags *LDFlags
	fs    afero.Fs
	logE  *logrus.Entry
}

func New(flags *LDFlags, fs afero.Fs, logE *logrus.Entry) *cli.Command {
	runner := &Runner{
		flags: flags,
		fs:    fs,
		logE:  logE,
	}
	return helpall.With(vcmd.With(&cli.Command{
		Name:                  "ghalint",
		Usage:                 "GitHub Actions linter",
		Version:               flags.Version,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
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
		},
		Commands: []*cli.Command{
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
		},
	}, flags.Commit), nil)
}
