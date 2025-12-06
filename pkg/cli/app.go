package cli

import (
	"context"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	fs := afero.NewOsFs()
	runner := &Runner{
		fs: fs,
	}
	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:  "ghalint",
		Usage: "GitHub Actions linter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-color",
				Usage: "log color",
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
				Action: urfave.Action(runner.Run, logger),
				Flags:  []cli.Flag{},
			},
			{
				Name: "run-action",
				Aliases: []string{
					"act",
				},
				Usage:  "lint actions",
				Action: urfave.Action(runner.RunAction, logger),
				Flags:  []cli.Flag{},
			},
			experiment.New(logger, fs),
		},
	}).Run(ctx, env.Args)
}

type Runner struct {
	fs afero.Fs
}
