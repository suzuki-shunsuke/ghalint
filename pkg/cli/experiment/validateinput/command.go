package validateinput

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller/schema"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/log"
	"github.com/urfave/cli/v3"
)

func New(logE *logrus.Entry, fs afero.Fs) *cli.Command {
	runner := &Runner{
		logE: logE,
		fs:   fs,
	}
	return &cli.Command{
		Name:        "validate-input",
		Usage:       "validate action inputs",
		Description: "validate action inputs",
		Action:      runner.Action,
	}
}

type Runner struct {
	logE *logrus.Entry
	fs   afero.Fs
}

func (r *Runner) Action(ctx context.Context, cmd *cli.Command) error {
	logE := r.logE

	if err := log.Set(logE, cmd.String("log-level"), cmd.String("log-color")); err != nil {
		return fmt.Errorf("configure log options: %w", err)
	}

	rootDir, err := GetRootDir()
	if err != nil {
		return fmt.Errorf("get the root directory: %w", err)
	}

	gh := github.New(ctx, logE)

	ctrl := schema.New(r.fs, logE, gh.Repositories, rootDir)

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
