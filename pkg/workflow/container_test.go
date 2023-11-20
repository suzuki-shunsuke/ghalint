package workflow_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestContainer_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		yaml  string
		image string
	}{
		{
			name:  "normal",
			yaml:  "image: node:18",
			image: "node:18",
		},
		{
			name:  "string",
			yaml:  "node:18",
			image: "node:18",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			c := &workflow.Container{}
			if err := yaml.Unmarshal([]byte(d.yaml), c); err != nil {
				t.Fatal(err)
			}
			if d.image != c.Image {
				t.Fatalf("got %v, wanted %v", c.Image, d.image)
			}
		})
	}
}
