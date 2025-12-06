package validateinput

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

func New(logger *slogutil.Logger, fs afero.Fs) *cli.Command {
	runner := &Runner{
		fs: fs,
	}
	return &cli.Command{
		Name:        "validate-input",
		Usage:       "validate action inputs",
		Description: "validate action inputs",
		Action:      urfave.Action(runner.Action, logger),
	}
}

type Runner struct {
	fs afero.Fs
}

func (r *Runner) Action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if err := logger.SetColor(cmd.String("log-color")); err != nil {
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
