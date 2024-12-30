package controller

import (
	"context"
	"encoding/json"
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
		&policy.JobTimeoutMinutesIsRequiredPolicy{},
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
		&policy.CheckoutPersistCredentialShouldBeFalsePolicy{},
	}
	var errs []*policy.ErrorInfo
	for _, filePath := range filePaths {
		logE := logE.WithField("workflow_file_path", filePath)
		for _, err := range c.validateWorkflow(logE, cfg, wfPolicies, jobPolicies, stepPolicies, filePath) {
			err.FilePath = filePath
			errs = append(errs, err)
		}
	}
	for _, err := range errs {
		if err.Policy != nil {
			err.Policy.URL = policy.GetURL(err.Policy.ID)
		}
	}
	if len(errs) > 0 {
		logE.Error("some workflow files are invalid")
		encoder := json.NewEncoder(c.stderr)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(errs); err != nil {
			logerr.WithError(logE, err).Error("encode the error info")
		}
		return debugError(errors.New("some workflow files are invalid"))
	}
	return nil
}

/*
policy name
file path
job name
step name, id, index, action name, version

file
location
policy name
policy description
policy url
error message
*/

func (c *Controller) validateWorkflow(logE *logrus.Entry, cfg *config.Config, wfPolicies []WorkflowPolicy, jobPolicies []JobPolicy, stepPolicies []StepPolicy, filePath string) []*policy.ErrorInfo {
	wf := &workflow.Workflow{
		FilePath: filePath,
	}
	if err := workflow.Read(c.fs, filePath, wf); err != nil {
		return []*policy.ErrorInfo{
			{
				Message: "read a workflow file: " + err.Error(),
			},
		}
	}

	wfCtx := &policy.WorkflowContext{
		FilePath: filePath,
		Workflow: wf,
	}

	var errs []*policy.ErrorInfo
	for _, wfPolicy := range wfPolicies {
		logE := withPolicyReference(logE, wfPolicy)
		if err := wfPolicy.ApplyWorkflow(logE, cfg, wfCtx, wf); err != nil {
			errs = append(errs, &policy.ErrorInfo{
				Policy: &policy.Info{
					Name:    wfPolicy.Name(),
					ID:      wfPolicy.ID(),
					Message: err.Error(),
				},
			})
		}
	}

	errs = append(errs, c.applyJobPolicies(logE, cfg, wfCtx, jobPolicies)...)
	errs = append(errs, c.applyStepPolicies(logE, cfg, wfCtx, wf.Jobs, stepPolicies)...)

	return errs
}

type Policy interface {
	Name() string
	ID() string
}

func withPolicyReference(logE *logrus.Entry, p Policy) *logrus.Entry {
	return logE.WithFields(logrus.Fields{
		"policy_name": p.Name(),
		"reference":   policy.GetURL(p.ID()),
	})
}

func (c *Controller) applyJobPolicies(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicies []JobPolicy) []*policy.ErrorInfo {
	var errs []*policy.ErrorInfo
	for _, jobPolicy := range jobPolicies {
		logE := withPolicyReference(logE, jobPolicy)
		for _, err := range c.applyJobPolicy(logE, cfg, wfCtx, jobPolicy) {
			err.Policy = &policy.Info{
				Name: jobPolicy.Name(),
				ID:   jobPolicy.ID(),
			}
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *Controller) applyJobPolicy(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobPolicy JobPolicy) []*policy.ErrorInfo {
	var errs []*policy.ErrorInfo
	for jobName, job := range wfCtx.Workflow.Jobs {
		jobCtx := &policy.JobContext{
			Workflow: wfCtx,
			Name:     jobName,
		}
		logE := logE.WithField("job_name", jobName)
		if err := jobPolicy.ApplyJob(logE, cfg, jobCtx, job); err != nil {
			errs = append(errs, &policy.ErrorInfo{
				Message: err.Error(),
				JobName: jobName,
			})
		}
	}
	return errs
}

func (c *Controller) applyStepPolicies(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicies []StepPolicy) []*policy.ErrorInfo {
	var errs []*policy.ErrorInfo
	for _, stepPolicy := range stepPolicies {
		logE := withPolicyReference(logE, stepPolicy)
		for _, err := range c.applyStepPolicy(logE, cfg, wfCtx, jobs, stepPolicy) {
			err.Policy = &policy.Info{
				Name: stepPolicy.Name(),
				ID:   stepPolicy.ID(),
			}
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *Controller) applyStepPolicy(logE *logrus.Entry, cfg *config.Config, wfCtx *policy.WorkflowContext, jobs map[string]*workflow.Job, stepPolicy StepPolicy) []*policy.ErrorInfo {
	var errs []*policy.ErrorInfo
	for jobName, job := range jobs {
		stepCtx := &policy.StepContext{
			FilePath: wfCtx.FilePath,
			Job: &policy.JobContext{
				Name:     jobName,
				Workflow: wfCtx,
				Job:      job,
			},
		}
		logE := logE.WithField("job_name", jobName)
		for idx, step := range job.Steps {
			logE := logE
			if step.ID != "" {
				logE = logE.WithField("step_id", step.ID)
			}
			if step.Name != "" {
				logE = logE.WithField("step_name", step.Name)
			}
			if err := stepPolicy.ApplyStep(logE, cfg, stepCtx, step); err != nil {
				errs = append(errs, &policy.ErrorInfo{
					StepName:  step.Name,
					StepID:    step.ID,
					StepIndex: idx,
					JobName:   jobName,
					Message:   err.Error(),
				})
			}
		}
	}
	return errs
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
