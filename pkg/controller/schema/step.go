package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"gopkg.in/yaml.v3"
)

type validateStep struct {
	step    *workflow.Step
	logger  *slog.Logger
	fs      afero.Fs
	gh      GitHub
	rootDir string
}

var fullCommitSHAPattern = regexp.MustCompile(`\b[0-9a-f]{40}\b`)

func (v *validateStep) readAction(ctx context.Context, action *workflow.Action) error { //nolint:cyclop
	if strings.HasPrefix(v.step.Uses, "./") {
		// local action
		if err := v.readLocalAction(action); err != nil {
			return fmt.Errorf("read a local action file: %w", err)
		}
		return nil
	}
	// <owner>/<repo>[/<path>]@<ref>
	fullPath, ref, ok := strings.Cut(v.step.Uses, "@")
	if !ok {
		return fmt.Errorf("invalid action format: %s", v.step.Uses)
	}
	elems := strings.Split(fullPath, "/")
	owner := elems[0]
	repo := elems[1]
	path := strings.Join(elems[2:], "/")
	sha := ref
	if !fullCommitSHAPattern.MatchString(ref) {
		// Get SHA of actions
		s, _, err := v.gh.GetCommitSHA1(ctx, owner, repo, ref, "")
		if err != nil {
			return fmt.Errorf("get commit SHA1: %w", err)
		}
		sha = s
	}
	// Download actions and store them in $GHALINT_ROOT_DIR/actions
	// Check if the action file exists
	cachePath := filepath.Join(v.rootDir, "actions", owner, repo, sha, path, "action.yaml")
	if f, err := afero.Exists(v.fs, cachePath); err != nil {
		return fmt.Errorf("check if the action file exists: %w", err)
	} else if f {
		if err := workflow.ReadAction(v.fs, cachePath, action); err != nil {
			return fmt.Errorf("read a cached action file: %w", err)
		}
		return nil
	}
	// Download action.yaml or action.yml
	content, err := v.download(ctx, &downloadInput{
		Owner: owner,
		Repo:  repo,
		Path:  path,
		Ref:   sha,
	})
	if err != nil {
		return fmt.Errorf("download action file: %w", err)
	}
	// write action.yaml to $GHALINT_ROOT_DIR/actions/<owner>/<repo>/<path>
	if err := v.fs.MkdirAll(filepath.Dir(cachePath), dirPermission); err != nil {
		return fmt.Errorf("create action directory: %w", err)
	}
	if err := afero.WriteFile(v.fs, cachePath, []byte(content), filePermission); err != nil {
		return fmt.Errorf("write action file: %w", err)
	}
	if err := yaml.Unmarshal([]byte(content), action); err != nil {
		return fmt.Errorf("unmarshal action file: %w", err)
	}
	return nil
}

const (
	filePermission = 0o644
	dirPermission  = 0o755
)

type downloadInput struct {
	Owner string
	Repo  string
	Path  string
	Ref   string
}

func (v *validateStep) download(ctx context.Context, input *downloadInput) (string, error) {
	for _, file := range []string{"action.yaml", "action.yml"} {
		content, _, _, err := v.gh.GetContents(ctx, input.Owner, input.Repo, filepath.Join(input.Path, file), &github.RepositoryContentGetOptions{
			Ref: input.Ref,
		})
		if err != nil {
			slogerr.WithError(v.logger, err).Debug("get action file")
			continue
		}
		s, err := content.GetContent()
		if err != nil {
			return "", fmt.Errorf("get content: %w", err)
		}
		return s, nil
	}
	return "", errors.New("action file can't be downloaded")
}

func (v *validateStep) validate(ctx context.Context) error {
	// Validate inputs
	if v.step.Uses == "" {
		return nil
	}
	v.logger = v.logger.With("action", v.step.Uses)
	action := &workflow.Action{}
	if err := v.readAction(ctx, action); err != nil {
		return fmt.Errorf("read action: %w", err)
	}
	validKeys := map[string]struct{}{}
	requiredKeys := map[string]struct{}{}
	validKeysArray := make([]string, 0, len(action.Inputs))
	requiredKeysArray := []string{}
	for key, input := range action.Inputs {
		validKeysArray = append(validKeysArray, key)
		validKeys[key] = struct{}{}
		if input.Required {
			requiredKeys[key] = struct{}{}
			requiredKeysArray = append(requiredKeysArray, key)
		}
	}
	validKeysS := strings.Join(validKeysArray, ", ")
	requiredKeysS := strings.Join(requiredKeysArray, ", ")
	v.logger = v.logger.With(
		"valid_inputs", validKeysS,
		"required_inputs", requiredKeysS,
	)
	failed := false
	// Check if the input is valid
	for key := range v.step.With {
		if _, ok := action.Inputs[key]; !ok {
			v.logger.Error("invalid input key", "input_key", key)
			failed = true
		}
	}
	// Check if required keys are set
	for key := range requiredKeys {
		if _, ok := v.step.With[key]; !ok {
			v.logger.Error("required key is not set", "input_key", key)
			failed = true
		}
	}
	if failed {
		return ErrSilent
	}
	return nil
}

func (v *validateStep) readLocalAction(action *workflow.Action) error {
	found := false
	for _, file := range []string{"action.yaml", "action.yml"} {
		p := filepath.Join(v.step.Uses, file)
		if f, err := afero.Exists(v.fs, p); err != nil {
			return fmt.Errorf("check if the action file exists: %w", err)
		} else if !f {
			continue
		}
		found = true
		if err := workflow.ReadAction(v.fs, p, action); err != nil {
			return fmt.Errorf("read a local action file: %w", err)
		}
	}
	if !found {
		return errors.New("local action file not found")
	}
	return nil
}
