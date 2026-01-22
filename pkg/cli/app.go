package cli

import (
	"context"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/experiment/validateinput"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/gflags"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type RunArgs struct {
	*gflags.GlobalFlags
}

type RunActionArgs struct {
	*gflags.GlobalFlags

	Files []string
}

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error { //nolint:funlen
	fs := afero.NewOsFs()
	runner := &Runner{
		fs: fs,
	}
	gf := &gflags.GlobalFlags{}
	runArgs := &RunArgs{
		GlobalFlags: gf,
	}
	runActionArgs := &RunActionArgs{
		GlobalFlags: gf,
	}
	validateInputArgs := &validateinput.Args{
		GlobalFlags: gf,
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
				Destination: &gf.LogColor,
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level",
				Sources: cli.EnvVars(
					"GHALINT_LOG_LEVEL",
				),
				Destination: &gf.LogLevel,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configuration file path",
				Sources: cli.EnvVars(
					"GHALINT_CONFIG",
				),
				Destination: &gf.Config,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "lint GitHub Actions Workflows",
				Action: func(ctx context.Context, _ *cli.Command) error {
					return runner.Run(ctx, logger, runArgs)
				},
			},
			{
				Name: "run-action",
				Aliases: []string{
					"act",
				},
				Usage: "lint actions",
				Action: func(ctx context.Context, _ *cli.Command) error {
					return runner.RunAction(ctx, logger, runActionArgs)
				},
				Arguments: []cli.Argument{
					&cli.StringArgs{
						Name:        "files",
						Destination: &runActionArgs.Files,
						Max:         -1,
					},
				},
			},
			experiment.New(logger, fs, validateInputArgs),
		},
	}).Run(ctx, env.Args)
}

type Runner struct {
	fs afero.Fs
}
