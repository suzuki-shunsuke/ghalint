package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

func (c *Controller) Run(_ context.Context, logger *slog.Logger, cfgFilePath string) error {
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
		&policy.JobTimeoutMinutesIsRequiredPolicy{},
		policy.NewJobSecretsPolicy(),
		&policy.DenyInheritSecretsPolicy{},
		&policy.DenyJobContainerLatestImagePolicy{},
		policy.NewActionRefShouldBeSHAPolicy(),
		&policy.DenyReadAllPermissionPolicy{},
		&policy.DenyWriteAllPermissionPolicy{},
	}
	stepPolicies := []StepPolicy{
		&policy.GitHubAppShouldLimitRepositoriesPolicy{},
		&policy.GitHubAppShouldLimitPermissionsPolicy{},
		policy.NewActionRefShouldBeSHAPolicy(),
		&policy.CheckoutPersistCredentialShouldBeFalsePolicy{},
	}
	failed := false
	for _, filePath := range filePaths {
		logger := logger.With("workflow_file_path", filePath)
		if c.validateWorkflow(logger, cfg, wfPolicies, jobPolicies, stepPolicies, filePath) {
			failed = true
		}
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}

func (c *Controller) validateWorkflow(logger *slog.Logger, cfg *config.Config, wfPolicies []WorkflowPolicy, jobPolicies []JobPolicy, stepPolicies []StepPolicy, filePath string) bool {
	wf := &workflow.Workflow{
		FilePath: filePath,
	}
	if err := workflow.Read(c.fs, filePath, wf); err != nil {
		slogerr.WithError(logger, err).Error("read a workflow file")
		return true
	}

	wfCtx := &policy.WorkflowContext{
		FilePath: filePath,
		Workflow: wf,
	}

	failed := false
	for _, wfPolicy := range wfPolicies {
		logger := withPolicyReference(logger, wfPolicy)
		if err := wfPolicy.ApplyWorkflow(logger, cfg, wfCtx, wf); err != nil {
			if err.Error() != "" {
				slogerr.WithError(logger, err).Error("the workflow violates policies")
			}
			failed = true
			continue
		}
	}

	if c.applyJobPolicies(logger, cfg, wfCtx, jobPolicies) {
		failed = true
	}

	if c.applyStepPolicies(logger, cfg, wfCtx, wf.Jobs, stepPolicies) {
		failed = true
	}

	return failed
}

type Policy interface {
	Name() string
	ID() string
}

func withPolicyReference(logger *slog.Logger, p Policy) *slog.Logger {
	return logger.With(
		"policy_name", p.Name(),
		"reference", fmt.Sprintf("https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/%s.md", p.ID()),
	)
}

func (c *Controller) applyJobPolicies(logger *slog.Logger, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicies []JobPolicy) bool {
	failed := false
	for _, jobPolicy := range jobPolicies {
		logger := withPolicyReference(logger, jobPolicy)
		if c.applyJobPolicy(logger, cfg, wfCtx, jobPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyJobPolicy(logger *slog.Logger, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicy JobPolicy) bool {
	failed := false
	for jobName, job := range wfCtx.Workflow.Jobs {
		jobCtx := &policy.JobContext{
			Workflow: wfCtx,
			Name:     jobName,
		}
		logger := logger.With("job_name", jobName)
		if err := jobPolicy.ApplyJob(logger, cfg, jobCtx, job); err != nil {
			failed = true
			if err.Error() != "" {
				slogerr.WithError(logger, err).Error("the job violates policies")
			}
		}
	}
	return failed
}

func (c *Controller) applyStepPolicies(logger *slog.Logger, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicies []StepPolicy) bool {
	failed := false
	for _, stepPolicy := range stepPolicies {
		logger := withPolicyReference(logger, stepPolicy)
		if c.applyStepPolicy(logger, cfg, wfCtx, jobs, stepPolicy) {
			failed = true
		}
	}
	return failed
}

func (c *Controller) applyStepPolicy(logger *slog.Logger, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicy StepPolicy) bool {
	failed := false
	for jobName, job := range jobs {
		stepCtx := &policy.StepContext{
			FilePath: wfCtx.FilePath,
			Job: &policy.JobContext{
				Name:     jobName,
				Workflow: wfCtx,
				Job:      job,
			},
		}
		logger := logger.With("job_name", jobName)
		for _, step := range job.Steps {
			logger := logger
			if step.ID != "" {
				logger = logger.With("step_id", step.ID)
			}
			if step.Name != "" {
				logger = logger.With("step_name", step.Name)
			}
			if err := stepPolicy.ApplyStep(logger, cfg, stepCtx, step); err != nil {
				if err.Error() != "" {
					slogerr.WithError(logger, err).Error("the step violates policies")
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
			return fmt.Errorf("read a configuration file: %w", slogerr.With(err,
				"config_file", cfgFilePath,
			))
		}
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validate a configuration file: %w", slogerr.With(err,
				"config_file", cfgFilePath,
			))
		}
		config.ConvertPath(cfg)
	}
	return nil
}
