package workflow

import (
	"errors"
)

type Container struct {
	Image string
}

func (c *Container) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val interface{}
	if err := unmarshal(&val); err != nil {
		return err
	}
	return convContainer(val, c)
}

func convContainer(src interface{}, c *Container) error { //nolint:cyclop
	switch p := src.(type) {
	case string:
		c.Image = p
		return nil
	case map[interface{}]interface{}:
		for k, v := range p {
			key, ok := k.(string)
			if !ok {
				continue
			}
			if key != "image" {
				continue
			}
			image, ok := v.(string)
			if !ok {
				return errors.New("image must be a string")
			}
			c.Image = image
			return nil
		}
		return nil
	case map[string]interface{}:
		for k, v := range p {
			if k != "image" {
				continue
			}
			image, ok := v.(string)
			if !ok {
				return errors.New("image must be a string")
			}
			c.Image = image
			return nil
		}
		return nil
	default:
		return errors.New("container must be a map or string")
	}
}
