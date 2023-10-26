package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Excludes []*Exclude
}

type Exclude struct {
	PolicyName       string `yaml:"policy_name"`
	WorkflowFilePath string `yaml:"workflow_file_path"`
	JobName          string `yaml:"job_name"`
}

func findConfig() string {
	for _, filePath := range []string{"ghalint.yaml", ".ghalint.yaml", "ghalint.yml", ".ghalint.yml"} {
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}
	return ""
}

func readConfig(cfg *Config, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open a configuration file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("parse configuration file as YAML: %w", err)
	}
	return nil
}

func validateConfig(cfg *Config) error {
	for _, exclude := range cfg.Excludes {
		if exclude.PolicyName == "" {
			return errors.New(`policy_name is required`)
		}
		if exclude.PolicyName != "job_secrets" {
			return errors.New(`only the policy "job_secrets" can be excluded`)
		}
		if exclude.WorkflowFilePath == "" {
			return errors.New(`workflow_file_path is required`)
		}
		if exclude.JobName == "" {
			return errors.New(`jobName is required`)
		}
	}
	return nil
}
