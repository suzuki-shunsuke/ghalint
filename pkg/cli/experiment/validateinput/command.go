package validateinput

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

func New(logger *slog.Logger, logLevelVar *slog.LevelVar, fs afero.Fs) *cli.Command {
	runner := &Runner{
		logger:      logger,
		logLevelVar: logLevelVar,
		fs:          fs,
	}
	return &cli.Command{
		Name:        "validate-input",
		Usage:       "validate action inputs",
		Description: "validate action inputs",
		Action:      runner.Action,
	}
}

type Runner struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
	fs          afero.Fs
}

func (r *Runner) Action(ctx context.Context, cmd *cli.Command) error {
	if cmd.String("log-color") != "" {
		r.logger.Warn("log color option is deprecated and doesn't work anymore. This is kept for backward compatibility.")
	}
	if err := slogutil.SetLevel(r.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	rootDir, err := GetRootDir()
	if err != nil {
		return fmt.Errorf("get the root directory: %w", err)
	}

	gh := github.New(ctx, r.logger)

	ctrl := schema.New(r.fs, r.logger, gh.Repositories, rootDir)

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
