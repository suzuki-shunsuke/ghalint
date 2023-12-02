package policy_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestGitHubAppShouldLimitRepositoriesPolicy_Apply(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		wf    *workflow.Workflow
		isErr bool
	}{
		{
			name:  "tibdex/github-app-token fail",
			isErr: true,
			cfg:   &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
								Uses: "tibdex/github-app-token@v2",
								ID:   "token",
								With: map[string]string{
									"app_id":      "xxx",
									"private_key": "xxx",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "tibdex/github-app-token success",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
								Uses: "tibdex/github-app-token@v2",
								ID:   "token",
								With: map[string]string{
									"app_id":       "xxx",
									"private_key":  "xxx",
									"repositories": "{}",
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "actions/create-github-app-token fail",
			isErr: true,
			cfg:   &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
								Uses: "actions/create-github-app-token@v2",
								ID:   "token",
								With: map[string]string{
									"app-id":      "xxx",
									"private-key": "xxx",
									"owner":       "xxx",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "actions/create-github-app-token success",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
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
					},
				},
			},
		},
		{
			name: "actions/create-github-app-token success no owner",
			cfg:  &config.Config{},
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
								Uses: "actions/create-github-app-token@v2",
								ID:   "token",
								With: map[string]string{
									"app-id":      "xxx",
									"private-key": "xxx",
								},
							},
						},
					},
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
			wf: &workflow.Workflow{
				FilePath: ".github/workflows/test.yaml",
				Jobs: map[string]*workflow.Job{
					"test": {
						Steps: []*workflow.Step{
							{
								Uses: "tibdex/github-app-token@v2",
								ID:   "token",
								With: map[string]string{
									"app_id":      "xxx",
									"private_key": "xxx",
								},
							},
						},
					},
				},
			},
		},
	}
	p := &policy.GitHubAppShouldLimitRepositoriesPolicy{}
	ctx := context.Background()
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.Apply(ctx, logE, d.cfg, d.wf); err != nil {
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
