package policy

import (
	"fmt"
	"os"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type WorkflowSecretsPolicy struct {
	secretPattern      *regexp.Regexp
	githubTokenPattern *regexp.Regexp
}

func NewWorkflowSecretsPolicy() *WorkflowSecretsPolicy {
	return &WorkflowSecretsPolicy{
		secretPattern:      regexp.MustCompile(`\${{ *secrets\.[^ ]+ *}}`),
		githubTokenPattern: regexp.MustCompile(`\${{ *github\.token+ *}}`),
	}
}

func (p *WorkflowSecretsPolicy) Name() string {
	return "workflow_secrets"
}

func (p *WorkflowSecretsPolicy) ID() string {
	return "005"
}

func (p *WorkflowSecretsPolicy) ApplyWorkflow(logE *logrus.Entry, _ *config.Config, wfctx *WorkflowContext, wf *workflow.Workflow) error {
	if len(wf.Jobs) < 2 { //nolint:mnd
		return nil
	}
	failed := false
	for envName, envValue := range wf.Env {
		logE := logE.WithField("env_name", envName)
		if p.secretPattern.MatchString(envValue) {
			failed = true
			logE.Error("secret should not be set to workflow's env")
			path, err := yaml.PathString("$.env." + envName)
			if err != nil {
				logerr.WithError(logE, err).Error("failed to create yaml path for env")
				continue
			}
			source, err := path.AnnotateSource(wfctx.Content, true)
			if err != nil {
				logerr.WithError(logE, err).Error("annotate yaml")
				continue
			}
			fmt.Fprintf(os.Stderr, "secret should not be set to workflow's env\n%s\n%s\n", wfctx.FilePath, string(source))
		}
		if p.githubTokenPattern.MatchString(envValue) {
			failed = true
			logE.WithField("env_name", envName).Error("github.token should not be set to workflow's env")
		}
	}
	if failed {
		return errEmpty
	}
	return nil
}
