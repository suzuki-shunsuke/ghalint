package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type DenyJobContainerLatestImagePolicy struct{}

func (p *DenyJobContainerLatestImagePolicy) Name() string {
	return "deny_job_container_latest_image"
}

func (p *DenyJobContainerLatestImagePolicy) Apply(ctx context.Context, logE *logrus.Entry, cfg *Config, wf *Workflow) error {
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
		return errors.New("workflow violates policies")
	}
	return nil
}
