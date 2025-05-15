package schema

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) runWorkflow(ctx context.Context) error {
	filePaths, err := workflow.List(c.fs)
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	failed := false
	for _, filePath := range filePaths {
		logE := c.logE.WithField("workflow_file_path", filePath)
		vw := &validateWorkflow{
			workflow: filePath,
			logE:     logE,
			fs:       c.fs,
			gh:       c.gh,
			rootDir:  c.rootDir,
		}
		if err := vw.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, SilentError) {
				logE.WithError(err).Error("validate workflow")
			}
		}
	}
	if failed {
		return SilentError
	}
	return nil
}

type validateWorkflow struct {
	workflow string
	logE     *logrus.Entry
	fs       afero.Fs
	gh       GitHub
	rootDir  string
}

func (v *validateWorkflow) validate(ctx context.Context) error {
	wf := &workflow.Workflow{
		FilePath: v.workflow,
	}
	if err := workflow.Read(v.fs, v.workflow, wf); err != nil {
		return fmt.Errorf("read a workflow file: %w", err)
	}
	failed := false
	for name, job := range wf.Jobs {
		vj := &validateJob{
			job:     job,
			logE:    v.logE.WithField("job_key", name),
			fs:      v.fs,
			gh:      v.gh,
			rootDir: v.rootDir,
		}
		if err := vj.validate(ctx); err != nil {
			failed = true
			if !errors.Is(err, SilentError) {
				logerr.WithError(v.logE, err).Error("validate job")
			}
		}
	}
	if failed {
		return SilentError
	}
	return nil
}
