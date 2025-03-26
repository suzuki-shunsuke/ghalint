package workflow

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

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
	TimeoutMinutes any `yaml:"timeout-minutes"`
}

type Step struct {
	Uses           string
	ID             string
	Name           string
	Run            string
	Shell          string
	With           With
	TimeoutMinutes any `yaml:"timeout-minutes"`
}

type With map[string]string

func (w With) UnmarshalYAML(b []byte) error {
	a := map[string]any{}
	if err := yaml.Unmarshal(b, &a); err != nil {
		return err //nolint:wrapcheck
	}
	for k, v := range a {
		switch c := v.(type) {
		case string:
			w[k] = c
		case int:
			w[k] = strconv.Itoa(c)
		case float64:
			w[k] = fmt.Sprint(c)
		case bool:
			w[k] = strconv.FormatBool(c)
		default:
			return fmt.Errorf("unsupported type: %T", c)
		}
	}
	return nil
}

type Action struct {
	Runs *Runs
}

type Runs struct {
	Image string
	Steps []*Step
}
