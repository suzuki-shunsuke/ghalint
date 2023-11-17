package controller

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type JobSecrets struct {
	m       map[string]string
	inherit bool
}

func (js *JobSecrets) Inherit() bool {
	return js != nil && js.inherit
}

func (js *JobSecrets) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val interface{}
	if err := unmarshal(&val); err != nil {
		return err
	}
	return convJobSecrets(val, js)
}

func convJobSecrets(src interface{}, dest *JobSecrets) error { //nolint:cyclop
	switch p := src.(type) {
	case string:
		switch p {
		case "inherit":
			dest.inherit = true
			return nil
		default:
			return logerr.WithFields(errors.New("job secrets must be a map or `inherit`"), logrus.Fields{ //nolint:wrapcheck
				"secrets": p,
			})
		}
	case map[interface{}]interface{}:
		m := make(map[string]string, len(p))
		for k, v := range p {
			ks, ok := k.(string)
			if !ok {
				return errors.New("secrets key must be string")
			}
			vs, ok := v.(string)
			if !ok {
				return errors.New("secrets value must be string")
			}
			m[ks] = vs
		}
		dest.m = m
		return nil
	case map[string]interface{}:
		m := make(map[string]string, len(p))
		for k, v := range p {
			vs, ok := v.(string)
			if !ok {
				return errors.New("secrets value must be string")
			}
			m[k] = vs
		}
		dest.m = m
		return nil
	default:
		return errors.New("secrets must be map[string]string or 'inherit'")
	}
}
