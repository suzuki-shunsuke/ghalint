package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestGitHubAppShouldLimitRepositoriesPolicy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		cfg     *config.Config
		stepCtx *policy.StepContext
		step    *workflow.Step
		isErr   bool
	}{
		{
			name:  "tibdex/github-app-token fail",
			isErr: true,
			cfg:   &config.Config{},
			step: &workflow.Step{
				Uses: "tibdex/github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app_id":      "xxx",
					"private_key": "xxx",
				},
			},
		},
		{
			name: "tibdex/github-app-token success",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "tibdex/github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app_id":       "xxx",
					"private_key":  "xxx",
					"repositories": "{}",
				},
			},
		},
		{
			name:  "actions/create-github-app-token fail",
			isErr: true,
			cfg:   &config.Config{},
			step: &workflow.Step{
				Uses: "actions/create-github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app-id":      "xxx",
					"private-key": "xxx",
					"owner":       "xxx",
				},
			},
		},
		{
			name: "actions/create-github-app-token success",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "actions/create-github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app-id":       "xxx",
					"private-key":  "xxx",
					"owner":        "xxx",
					"repositories": "foo,bar",
				},
			},
		},
		{
			name: "actions/create-github-app-token success no owner",
			cfg:  &config.Config{},
			step: &workflow.Step{
				Uses: "actions/create-github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app-id":      "xxx",
					"private-key": "xxx",
				},
			},
		},
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "github_app_should_limit_repositories",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "test",
						StepID:           "token",
					},
				},
			},
			stepCtx: &policy.StepContext{
				FilePath: ".github/workflows/test.yaml",
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yaml",
					},
				},
			},
			step: &workflow.Step{
				Uses: "tibdex/github-app-token@v2",
				ID:   "token",
				With: map[string]string{
					"app_id":      "xxx",
					"private_key": "xxx",
				},
			},
		},
	}
	p := &policy.GitHubAppShouldLimitRepositoriesPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		if d.stepCtx == nil {
			d.stepCtx = &policy.StepContext{
				FilePath: ".github/workflows/test.yaml",
				Job: &policy.JobContext{
					Name: "test",
					Workflow: &policy.WorkflowContext{
						FilePath: ".github/workflows/test.yaml",
					},
				},
			}
		}
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyStep(logE, d.cfg, d.stepCtx, d.step); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
