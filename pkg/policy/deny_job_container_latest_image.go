package policy

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

type DenyJobContainerLatestImagePolicy struct{}

func (p *DenyJobContainerLatestImagePolicy) Name() string {
	return "deny_job_container_latest_image"
}

func (p *DenyJobContainerLatestImagePolicy) ID() string {
	return "007"
}

func (p *DenyJobContainerLatestImagePolicy) Apply(_ context.Context, logE *logrus.Entry, _ *config.Config, wf *workflow.Workflow) error {
	failed := false
	for jobName, job := range wf.Jobs {
		logE := logE.WithField("job_name", jobName)
		if job.Container == nil {
			continue
		}
		if job.Container.Image == "" {
			logE.Error("job container should have image")
			failed = true
			continue
		}
		_, tag, ok := strings.Cut(job.Container.Image, ":")
		if !ok {
			logE.Error("job container image should be <image name>:<tag>")
			failed = true
			continue
		}
		if tag == "latest" {
			logE.Error("job container image tag should not be `latest`")
			failed = true
		}
	}
	if failed {
		return errWorkflowViolatePolicy
	}
	return nil
}
