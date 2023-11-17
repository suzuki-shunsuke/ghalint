package config_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
)

func TestValidate(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		cfg   *config.Config
		isErr bool
	}{
		{
			name: "policy_name is required",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{},
				},
			},
			isErr: true,
		},
		{
			name: "action_name is required",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "action_ref_should_be_full_length_commit_sha",
					},
				},
			},
			isErr: true,
		},
		{
			name: "workflow_file_path is required",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName: "job_secrets",
					},
				},
			},
			isErr: true,
		},
		{
			name: "job_name is required",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "job_secrets",
						WorkflowFilePath: ".github/workflows/foo.yaml",
					},
				},
			},
			isErr: true,
		},
		{
			name: "disallowed policy",
			cfg: &config.Config{
				Excludes: []*config.Exclude{
					{
						PolicyName:       "deny_read_all_permission",
						WorkflowFilePath: ".github/workflows/foo.yaml",
						JobName:          "foo",
					},
				},
			},
			isErr: true,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := config.Validate(d.cfg); err != nil {
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
