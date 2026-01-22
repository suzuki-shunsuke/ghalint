package validateinput

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/cli/gflags"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type Args struct {
	*gflags.GlobalFlags
}

func New(logger *slogutil.Logger, fs afero.Fs, args *Args) *cli.Command {
	runner := &Runner{
		fs: fs,
	}
	return &cli.Command{
		Name:        "validate-input",
		Usage:       "validate action inputs",
		Description: "validate action inputs",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runner.Action(ctx, logger, args)
		},
	}
}

type Runner struct {
	fs afero.Fs
}

func (r *Runner) Action(ctx context.Context, logger *slogutil.Logger, args *Args) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(args.LogColor); err != nil {
		return fmt.Errorf("set log color: %w", err)
	}

	rootDir, err := GetRootDir()
	if err != nil {
		return fmt.Errorf("get the root directory: %w", err)
	}

	gh := github.New(ctx, logger.Logger)

	ctrl := schema.New(r.fs, logger.Logger, gh.Repositories, rootDir)

	return ctrl.Run(ctx) //nolint:wrapcheck
}

func GetRootDir() (string, error) {
	// ${GHALINT_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/ghalint}
	rootDir := os.Getenv("GHALINT_ROOT_DIR")
	if rootDir != "" {
		return rootDir, nil
	}
	xdgDataHome := xdg.DataHome
	if xdgDataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get the current user home directory: %w", err)
		}
		xdgDataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdgDataHome, "ghalint"), nil
}
