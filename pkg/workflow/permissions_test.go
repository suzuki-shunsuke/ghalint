package workflow_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"gopkg.in/yaml.v3"
)

func TestPermissions_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		yaml     string
		readAll  bool
		writeAll bool
	}{
		{
			name: "not read-all and write-all",
			yaml: `contents: read`,
		},
		{
			name:    "read-all",
			yaml:    `read-all`,
			readAll: true,
		},
		{
			name:     "write-all",
			yaml:     `write-all`,
			writeAll: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			p := &workflow.Permissions{}
			if err := yaml.Unmarshal([]byte(d.yaml), p); err != nil {
				t.Fatal(err)
			}
			readAll := p.ReadAll()
			writeAll := p.WriteAll()
			if d.readAll != readAll {
				t.Fatalf("readAll got %v, wanted %v", readAll, d.readAll)
			}
			if d.writeAll != writeAll {
				t.Fatalf("writeAll got %v, wanted %v", writeAll, d.writeAll)
			}
		})
	}
}
