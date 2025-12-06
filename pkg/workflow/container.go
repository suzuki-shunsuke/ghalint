package workflow

import (
	"errors"
)

type Container struct {
	Image string
}

func (c *Container) UnmarshalYAML(unmarshal func(any) error) error {
	var val any
	if err := unmarshal(&val); err != nil {
		return err
	}
	return convContainer(val, c)
}

func convContainer(src any, c *Container) error { //nolint:cyclop
	switch p := src.(type) {
	case string:
		c.Image = p
		return nil
	case map[any]any:
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
	case map[string]any:
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
