package policy_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestJobSecrets_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		yaml    string
		inherit bool
	}{
		{
			name: "not inherit",
			yaml: `token: ${{github.token}}`,
		},
		{
			name:    "inherit",
			yaml:    `inherit`,
			inherit: true,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			js := &workflow.JobSecrets{}
			if err := yaml.Unmarshal([]byte(d.yaml), js); err != nil {
				t.Fatal(err)
			}
			inherit := js.Inherit()
			if d.inherit != inherit {
				t.Fatalf("got %v, wanted %v", inherit, d.inherit)
			}
		})
	}
}
