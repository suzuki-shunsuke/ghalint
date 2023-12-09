package policy

import (
	"errors"
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

func (p *DenyJobContainerLatestImagePolicy) ApplyJob(_ *logrus.Entry, _ *config.Config, _ *JobContext, job *workflow.Job) error {
	if job.Container == nil {
		return nil
	}
	if job.Container.Image == "" {
		return errors.New("job container should have image")
	}
	_, tag, ok := strings.Cut(job.Container.Image, ":")
	if !ok {
		return errors.New("job container image should be <image name>:<tag>")
	}
	if tag == "latest" {
		return errors.New("job container image tag should not be `latest`")
	}
	return nil
}
