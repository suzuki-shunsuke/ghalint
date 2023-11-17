package workflow

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Permissions struct {
	m        map[string]string
	readAll  bool
	writeAll bool
}

func (ps *Permissions) Permissions() map[string]string {
	if ps == nil {
		return nil
	}
	return ps.m
}

func (ps *Permissions) ReadAll() bool {
	if ps == nil {
		return false
	}
	return ps.readAll
}

func (ps *Permissions) WriteAll() bool {
	if ps == nil {
		return false
	}
	return ps.writeAll
}

func (ps *Permissions) IsNil() bool {
	if ps == nil {
		return true
	}
	return ps.m == nil && !ps.readAll && !ps.writeAll
}

func (ps *Permissions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val interface{}
	if err := unmarshal(&val); err != nil {
		return err
	}
	return convPermissions(val, ps)
}

func convPermissions(src interface{}, dest *Permissions) error { //nolint:cyclop
	switch p := src.(type) {
	case string:
		switch p {
		case "read-all":
			dest.readAll = true
			return nil
		case "write-all":
			dest.writeAll = true
			return nil
		default:
			return logerr.WithFields(errors.New("unknown permissions"), logrus.Fields{ //nolint:wrapcheck
				"permission": p,
			})
		}
	case map[interface{}]interface{}:
		m := make(map[string]string, len(p))
		for k, v := range p {
			ks, ok := k.(string)
			if !ok {
				return errors.New("permissions key must be string")
			}
			vs, ok := v.(string)
			if !ok {
				return errors.New("permissions value must be string")
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
				return errors.New("permissions value must be string")
			}
			m[k] = vs
		}
		dest.m = m
		return nil
	default:
		return errors.New("permissions must be map[string]string or 'read-all' or 'write-all'")
	}
}
