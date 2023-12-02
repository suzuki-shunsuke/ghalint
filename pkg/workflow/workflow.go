package workflow

type Workflow struct {
	FilePath    string `yaml:"-"`
	Jobs        map[string]*Job
	Env         map[string]string
	Permissions *Permissions
}

type Job struct {
	Permissions *Permissions
	Env         map[string]string
	Steps       []*Step
	Secrets     *JobSecrets
	Container   *Container
	Uses        string
}

type Step struct {
	Uses string
	ID   string
	Name string
	With map[string]string
}
