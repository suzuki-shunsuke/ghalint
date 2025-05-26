package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/action"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) runActions(ctx context.Context) error {
	filePaths, err := action.Find(c.fs)
	if err != nil {
		return fmt.Errorf("find action files: %w", err)
	}
	failed := false
	for _, filePath := range filePaths {
		logE := c.logE.WithField("action_file_path", filePath)
		vw := &validateAction{
			action:  filePath,
			logE:    logE,
			fs:      c.fs,
			gh:      c.gh,
			rootDir: c.rootDir,
		}
		if err := vw.validate(ctx); err != nil {
			logE.WithError(err).Error("validate action")
			failed = true
		}
	}
	if failed {
		return errors.New("some action files are invalid")
	}
	return nil
}

type validateAction struct {
	action  string
	logE    *logrus.Entry
	fs      afero.Fs
	gh      GitHub
	rootDir string
}

func (v *validateAction) validate(ctx context.Context) error {
	act := &workflow.Action{}
	if err := workflow.ReadAction(v.fs, v.action, act); err != nil {
		return fmt.Errorf("read an action file: %w", err)
	}
	failed := false
	for _, step := range act.Runs.Steps {
		vs := &validateStep{
			step:    step,
			logE:    v.logE,
			fs:      v.fs,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vs.validate(ctx); err != nil {
			logerr.WithError(v.logE, err).Error("validate a step")
			failed = true
		}
	}
	if failed {
		return errors.New("some steps are invalid")
	}
	return nil
}
