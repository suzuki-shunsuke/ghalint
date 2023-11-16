package cli

type Workflow struct {
	FilePath    string `yaml:"-"`
	Jobs        map[string]*Job
	Env         map[string]string
	Permissions *Permissions
}

type Job struct {
	Permissions *Permissions
	Env         map[string]string
	Steps       []interface{}
	Secrets     *JobSecrets
	Container   *Container
}

type Container struct {
	Image string
}
