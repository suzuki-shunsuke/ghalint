package cli

type Config struct {
	Excludes []*Exclude
}

type Exclude struct {
	PolicyName       string `yaml:"policy_name"`
	WorkflowFilePath string `yaml:"workflow_file_path"`
	JobName          string `yaml:"job_name"`
}
