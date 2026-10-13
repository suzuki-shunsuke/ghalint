package config_test

import (
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/config"
)

func TestMatchActionName(t *testing.T) {
	t.Parallel()
	cases := []struct {
		pattern string
		name    string
		want    bool
	}{
		{"my-org/actions/*", "my-org/actions/setup", true},
		{"my-org/actions/*", "my-org/actions/setup/node", false},
		{"my-org/actions/**", "my-org/actions/setup", true},
		{"my-org/actions/**", "my-org/actions/setup/node", true},
		{"my-org/actions/**", "my-org/actions/a/b/c", true},
		{"my-org/actions/**", "other/actions/setup", false},
		{"docker://rhysd/actionlint", "docker://rhysd/actionlint", true},
	}
	for _, tc := range cases {
		t.Run(tc.pattern+"::"+tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := config.MatchActionName(tc.pattern, tc.name)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
