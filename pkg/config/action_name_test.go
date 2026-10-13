package config_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
)

func TestMatchActionName(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		pattern string
		action  string
		want    bool
	}{
		{
			name:    "path match without double star",
			pattern: "suzuki-shunsuke/tfaction/*",
			action:  "suzuki-shunsuke/tfaction/pinact",
			want:    true,
		},
		{
			name:    "single star does not match slash",
			pattern: "suzuki-shunsuke/tfaction/*",
			action:  "suzuki-shunsuke/tfaction/nested/pinact",
		},
		{
			name:    "double star matches nested action name",
			pattern: "my-private-org/actions/**",
			action:  "my-private-org/actions/foo/bar",
			want:    true,
		},
		{
			name:    "double star matches zero action name parts",
			pattern: "my-private-org/actions/**",
			action:  "my-private-org/actions",
			want:    true,
		},
		{
			name:    "double star matches middle action name parts",
			pattern: "my-private-org/**/setup",
			action:  "my-private-org/actions/foo/setup",
			want:    true,
		},
		{
			name:    "double star respects other pattern parts",
			pattern: "my-private-org/actions/**",
			action:  "other-private-org/actions/foo/bar",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			got, err := config.MatchActionName(d.pattern, d.action)
			if err != nil {
				t.Fatal(err)
			}
			if got != d.want {
				t.Fatalf("wanted %t, got %t", d.want, got)
			}
		})
	}
}

func TestMatchActionName_InvalidPattern(t *testing.T) {
	t.Parallel()
	if _, err := config.MatchActionName("my-private-org/actions/[", "my-private-org/actions/foo"); err == nil {
		t.Fatal("error must be returned")
	}
}
