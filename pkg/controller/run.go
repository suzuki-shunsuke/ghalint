package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) Run(_ context.Context, logE *logrus.Entry, cfgFilePath string) error {
	cfg := &config.Config{}
	if err := c.readConfig(cfg, cfgFilePath); err != nil {
		return err
	}

	filePaths, err := workflow.List(c.fs)
	if err != nil {
		return fmt.Errorf("find workflow files: %w", err)
	}
	wfPolicies := []WorkflowPolicy{
		policy.NewWorkflowSecretsPolicy(),
	}
	jobPolicies := []JobPolicy{
		&policy.JobPermissionsPolicy{},
		policy.NewJobSecretsPolicy(),
		&policy.DenyInheritSecretsPolicy{},
		&policy.DenyJobContainerLatestImagePolicy{},
		policy.NewActionRefShouldBeSHA1Policy(),
		&policy.DenyReadAllPermissionPolicy{},
		&policy.DenyWriteAllPermissionPolicy{},
	}
	stepPolicies := []StepPolicy{
		&policy.GitHubAppShouldLimitRepositoriesPolicy{},
		&policy.GitHubAppShouldLimitPermissionsPolicy{},
		policy.NewActionRefShouldBeSHA1Policy(),
	}
	failed := false
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		if c.validateWorkflow(logE, cfg, wfPolicies, jobPolicies, stepPolicies, filePath) {
			failed = true
		}
	}
	if failed {
		return debugError(errors.New("some workflow files are invalid"))
	}
	return nil
}

func (c *Controller) validateWorkflow(logE *logrus.Entry, cfg *config.Config, wfPolicies []WorkflowPolicy, jobPolicies []JobPolicy, stepPolicies []StepPolicy, filePath string) bool {
	wf := &workflow.Workflow{
		FilePath: filePath,
	}
	if err := workflow.Read(c.fs, filePath, wf); err != nil {
		logerr.WithError(logE, err).Error("read a workflow file")
		return true
	}

	wfCtx := &policy.WorkflowContext{
		FilePath: filePath,
		Workflow: wf,
	}

	failed := false
	for _, wfPolicy := range wfPolicies {
		logE := withPolicyReference(logE, wfPolicy)
		if err := wfPolicy.ApplyWorkflow(logE, cfg, wfCtx, wf); err != nil {
			if err.Error() != "" {
				logerr.WithError(logE, err).Error("the workflow violates policies")
			}
			failed = true
			continue
		}
	}

	if c.applyJobPolicies(logE, cfg, wfCtx, jobPolicies) {
		failed = true
	}

	if c.applyStepPolicies(logE, cfg, wfCtx, wf.Jobs, stepPolicies) {
		failed = true
	}

	return failed
}

type Policy interface {
	Name() string
	ID() string
}

func withPolicyReference(logE *logrus.Entry, p Policy) *logrus.Entry {
	return logE.WithFields(logrus.Fields{
		"policy_name": p.Name(),
		"reference":   fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", p.ID()),
	})
}

func (c *Controller) applyJobPolicies(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicies []JobPolicy) bool {
	failed := false
	for _, jobPolicy := range jobPolicies {
		logE := withPolicyReference(logE, jobPolicy)
		if c.applyJobPolicy(logE, cfg, wfCtx, jobPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyJobPolicy(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicy JobPolicy) bool {
	failed := false
	for jobName, job := range wfCtx.Workflow.Jobs {
		jobCtx := &policy.JobContext{
			Workflow: wfCtx,
			Name:     jobName,
		}
		logE := logE.WithField("job_name", jobName)
		if err := jobPolicy.ApplyJob(logE, cfg, jobCtx, job); err != nil {
			failed = true
			if err.Error() != "" {
				logerr.WithError(logE, err).Error("the job violates policies")
			}
		}
	}
	return failed
}

func (c *Controller) applyStepPolicies(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicies []StepPolicy) bool {
	failed := false
	for _, stepPolicy := range stepPolicies {
		logE := withPolicyReference(logE, stepPolicy)
		if c.applyStepPolicy(logE, cfg, wfCtx, jobs, stepPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyStepPolicy(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicy StepPolicy) bool {
	failed := false
	for jobName, job := range jobs {
		jobCtx := &policy.JobContext{
			Workflow: wfCtx,
			Name:     jobName,
		}
		logE := logE.WithField("job_name", jobName)
		for _, step := range job.Steps {
			logE := logE
			if step.ID != "" {
				logE = logE.WithField("step_id", step.ID)
			}
			if step.Name != "" {
				logE = logE.WithField("step_name", step.Name)
			}
			if err := stepPolicy.ApplyStep(logE, cfg, jobCtx, step); err != nil {
				if err.Error() != "" {
					logerr.WithError(logE, err).Error("the step violates policies")
				}
				failed = true
			}
		}
	}
	return failed
}

func (c *Controller) readConfig(cfg *config.Config, cfgFilePath string) error {
	if cfgFilePath == "" {
		if c := config.Find(c.fs); c != "" {
			cfgFilePath = c
		}
	}
	if cfgFilePath != "" {
		if err := config.Read(c.fs, cfg, cfgFilePath); err != nil {
			return fmt.Errorf("read a configuration file: %w", logerr.WithFields(err, logrus.Fields{
				"config_file": cfgFilePath,
			}))
		}
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validate a configuration file: %w", logerr.WithFields(err, logrus.Fields{
				"config_file": cfgFilePath,
			}))
		}
	}
	return nil
}
