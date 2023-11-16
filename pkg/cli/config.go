package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Excludes []*Exclude
}

type Exclude struct {
	PolicyName       string `yaml:"policy_name"`
	WorkflowFilePath string `yaml:"workflow_file_path"`
	JobName          string `yaml:"job_name"`
	ActionName       string `yaml:"action_name"`
}

func findConfig(fs afero.Fs) string {
	for _, filePath := range []string{"ghalint.yaml", ".ghalint.yaml", "ghalint.yml", ".ghalint.yml"} {
		if _, err := fs.Stat(filePath); err == nil {
			return filePath
		}
	}
	return ""
}

func readConfig(fs afero.Fs, cfg *Config, filePath string) error {
	f, err := fs.Open(filePath)
	if err != nil {
		return fmt.Errorf("open a configuration file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("parse configuration file as YAML: %w", err)
	}
	return nil
}

func validateConfig(cfg *Config) error {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == "" {
			return errors.New(`policy_name is required`)
		}
		switch exclude.PolicyName {
		case "action_ref_should_be_full_length_commit_sha":
			if exclude.ActionName == "" {
				return errors.New(`action_name is required to exclude action_ref_should_be_full_length_commit_sha`)
			}
		case "job_secrets":
			if exclude.WorkflowFilePath == "" {
				return errors.New(`workflow_file_path is required`)
			}
			if exclude.JobName == "" {
				return errors.New(`jobName is required`)
			}
		default:
			return errors.New(`only the policy "job_secrets" and "action_ref_should_be_full_length_commit_sha" can be excluded`)
		}
	}
	return nil
}
