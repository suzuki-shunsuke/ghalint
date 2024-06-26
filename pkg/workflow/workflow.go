package workflow

type Workflow struct {
	FilePath    string `yaml:"-"`
	Jobs        map[string]*Job
	Env         map[string]string
	Permissions *Permissions
}

type Job struct {
	Permissions    *Permissions
	Env            map[string]string
	Steps          []*Step
	Secrets        *JobSecrets
	Container      *Container
	Uses           string
	TimeoutMinutes int `yaml:"timeout-minutes"`
}

type Step struct {
	Uses           string
	ID             string
	Name           string
	Run            string
	Shell          string
	With           map[string]string
	TimeoutMinutes int `yaml:"timeout-minutes"`
}

type Action struct {
	Runs *Runs
}

type Runs struct {
	Image string
	Steps []*Step
}
