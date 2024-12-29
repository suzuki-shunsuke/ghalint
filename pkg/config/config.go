package config

import (
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Excludes []*Exclude `json:"excludes,omitempty"`
}

type Exclude struct {
	PolicyName       string `json:"policy_name" yaml:"policy_name"`
	WorkflowFilePath string `json:"workflow_file_path,omitempty" yaml:"workflow_file_path"`
	ActionFilePath   string `json:"action_file_path,omitempty" yaml:"action_file_path"`
	JobName          string `json:"job_name,omitempty" yaml:"job_name"`
	ActionName       string `json:"action_name,omitempty" yaml:"action_name"`
	StepID           string `json:"step_id,omitempty" yaml:"step_id"`
}

func (e *Exclude) FilePath() string {
	if e.WorkflowFilePath != "" {
		return e.WorkflowFilePath
	}
	return e.ActionFilePath
}

func Find(fs afero.Fs) string {
	for _, filePath := range []string{"ghalint.yaml", ".ghalint.yaml", "ghalint.yml", ".ghalint.yml"} {
		if _, err := fs.Stat(filePath); err == nil {
			return filePath
		}
	}
	return ""
}

func Read(fs afero.Fs, cfg *Config, filePath string) error {
	f, err := fs.Open(filePath)
	if err != nil {
		return fmt.Errorf("open a configuration file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		err := fmt.Errorf("parse configuration file as YAML: %w", err)
		if errors.Is(err, io.EOF) {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"reference": "https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/codes/002.md",
			})
		}
		return err
	}
	return nil
}

func Validate(cfg *Config) error {
	for _, exclude := range cfg.Excludes {
		if err := validate(exclude); err != nil {
			return err
		}
	}
	return nil
}

func validate(exclude *Exclude) error { //nolint:cyclop
	if exclude.PolicyName == "" {
		return errors.New(`policy_name is required`)
	}
	switch exclude.PolicyName {
	case "action_ref_should_be_full_length_commit_sha":
		if exclude.ActionName == "" {
			return errors.New(`action_name is required to exclude action_ref_should_be_full_length_commit_sha`)
		}
		if _, err := path.Match(exclude.ActionName, ""); err != nil {
			return fmt.Errorf("action_name must be a glob pattern: %w", logerr.WithFields(err, logrus.Fields{
				"pattern_reference": "https://pkg.go.dev/path#Match",
			}))
		}
	case "job_secrets":
		if exclude.WorkflowFilePath == "" {
			return errors.New(`workflow_file_path is required to exclude job_secrets`)
		}
		if exclude.JobName == "" {
			return errors.New(`job_name is required to exclude job_secrets`)
		}
	case "deny_inherit_secrets":
		if exclude.WorkflowFilePath == "" {
			return errors.New(`workflow_file_path is required to exclude deny_inherit_secrets`)
		}
		if exclude.JobName == "" {
			return errors.New(`job_name is required to exclude deny_inherit_secrets`)
		}
	case "github_app_should_limit_repositories":
		if exclude.WorkflowFilePath == "" && exclude.ActionFilePath == "" {
			return errors.New(`workflow_file_path or action_file_path is required to exclude github_app_should_limit_repositories`)
		}
		if exclude.WorkflowFilePath != "" && exclude.JobName == "" {
			return errors.New(`job_name is required to exclude github_app_should_limit_repositories`)
		}
		if exclude.StepID == "" {
			return errors.New(`step_id is required to exclude github_app_should_limit_repositories`)
		}
	case "checkout_persist_credentials_should_be_false":
		if exclude.WorkflowFilePath == "" && exclude.ActionFilePath == "" {
			return errors.New(`workflow_file_path or action_file_path is required to exclude checkout_persist_credentials_should_be_false`)
		}
		if exclude.WorkflowFilePath != "" && exclude.JobName == "" {
			return errors.New(`job_name is required to exclude checkout_persist_credentials_should_be_false`)
		}
	default:
		return logerr.WithFields(errors.New(`the policy can't be excluded`), logrus.Fields{ //nolint:wrapcheck
			"policy_name": exclude.PolicyName,
		})
	}
	return nil
}
