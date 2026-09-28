package config

import (
	"fmt"
	"path"
	"strings"
)

// MatchActionName reports whether pattern matches action.
func MatchActionName(pattern, action string) (bool, error) {
	if !strings.Contains(pattern, "**") {
		matched, err := path.Match(pattern, action)
		if err != nil {
			return false, fmt.Errorf("match action name: %w", err)
		}
		return matched, nil
	}
	if err := validateActionNamePattern(pattern); err != nil {
		return false, err
	}
	return matchActionNameParts(strings.Split(pattern, "/"), strings.Split(action, "/"))
}

func validateActionNamePattern(pattern string) error {
	for part := range strings.SplitSeq(pattern, "/") {
		if part == "**" {
			continue
		}
		if _, err := path.Match(part, ""); err != nil {
			return fmt.Errorf("validate action name pattern: %w", err)
		}
	}
	return nil
}

func matchActionNameParts(patternParts, actionParts []string) (bool, error) {
	if len(patternParts) == 0 {
		return len(actionParts) == 0, nil
	}
	if patternParts[0] == "**" {
		for i := range len(actionParts) + 1 {
			matched, err := matchActionNameParts(patternParts[1:], actionParts[i:])
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	}
	if len(actionParts) == 0 {
		return false, nil
	}
	matched, err := path.Match(patternParts[0], actionParts[0])
	if err != nil {
		return false, fmt.Errorf("match action name part: %w", err)
	}
	if !matched {
		return false, nil
	}
	return matchActionNameParts(patternParts[1:], actionParts[1:])
}
