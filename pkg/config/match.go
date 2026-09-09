package config

import (
	"path"
	"regexp"
	"strings"
)

// MatchActionName matches an action name against an exclude pattern.
//
// Patterns use the same rules as path.Match, and additionally support "**"
// which matches across path separators (including zero segments). This lets
// excludes like "my-org/actions/**" cover nested composite actions.
func MatchActionName(pattern, name string) (bool, error) {
	if !strings.Contains(pattern, "**") {
		return path.Match(pattern, name)
	}
	re, err := doublestarToRegexp(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(name), nil
}

func doublestarToRegexp(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteByte('^')
	for i := 0; i < len(pattern); {
		if i+1 < len(pattern) && pattern[i] == '*' && pattern[i+1] == '*' {
			b.WriteString(".*")
			i += 2
			if i < len(pattern) && pattern[i] == '/' {
				// "**/" also absorbs the following slash so "a/**/b" works.
				i++
			}
			continue
		}
		switch c := pattern[i]; c {
		case '*':
			b.WriteString("[^/]*")
		case '?':
			b.WriteString("[^/]")
		case '.', '+', '(', ')', '|', '^', '$', '{', '}', '[', ']', '\':
			b.WriteByte('\\')
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
		i++
	}
	b.WriteByte('$')
	return regexp.Compile(b.String())
}
