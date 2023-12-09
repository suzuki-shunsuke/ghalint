package policy_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestGitHubAppShouldLimitPermissionsPolicy_ApplyStep(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		cfg    *config.Config
		jobCtx *policy.JobContext
		step   *workflow.Step
		isErr  bool
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
					"app_id":      "xxx",
					"private_key": "xxx",
					"permissions": "{}",
				},
			},
		},
		{
			name: "exclude",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "github_app_should_limit_permissions",
						WorkflowFilePath: ".github/workflows/test.yaml",
						JobName:          "test",
						StepID:           "token",
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
	p := &policy.GitHubAppShouldLimitPermissionsPolicy{}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		if d.jobCtx == nil {
			d.jobCtx = &policy.JobContext{
				Workflow: &policy.WorkflowContext{
					FilePath: ".github/workflows/test.yaml",
				},
				Name: "test",
			}
		}
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyStep(logE, d.cfg, d.jobCtx, d.step); err != nil {
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
