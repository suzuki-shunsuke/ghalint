package schema

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"gopkg.in/yaml.v3"
)

type validateStep struct {
	step    *workflow.Step
	logger  *slog.Logger
	gh      GitHub
	rootDir string
}

func (v *validateStep) readAction(ctx context.Context, action *workflow.Action) error { //nolint:cyclop
	if strings.HasPrefix(v.step.Uses, "./") {
		// local action
		if err := v.readLocalAction(action); err != nil {
			return fmt.Errorf("read a local action file: %w", err)
		}
		return nil
	}
	r, err := parseRemoteRef(v.step.Uses)
	if err != nil {
		return fmt.Errorf("parse the action: %w", err)
	}
	sha, err := r.resolveSHA(ctx, v.gh)
	if err != nil {
		return err
	}
	// Download actions and store them in $GHALINT_ROOT_DIR/actions
	// Check if the action file exists
	cachePath, err := r.cachePath(v.rootDir, sha, "action.yaml")
	if err != nil {
		return fmt.Errorf("get the cache path: %w", err)
	}
	if err := workflow.ReadAction(cachePath, action); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read a cached action file: %w", err)
	}
	// Download action.yaml or action.yml
	content, err := v.download(ctx, &downloadInput{
		Owner: r.Owner,
		Repo:  r.Repo,
		Path:  r.Path,
		Ref:   sha,
	})
	if err != nil {
		return fmt.Errorf("download action file: %w", err)
	}
	// write action.yaml to $GHALINT_ROOT_DIR/actions/<owner>/<repo>/<path>
	if err := os.MkdirAll(filepath.Dir(cachePath), dirPermission); err != nil {
		return fmt.Errorf("create action directory: %w", err)
	}
	if err := os.WriteFile(cachePath, []byte(content), filePermission); err != nil {
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
		content, _, _, err := v.gh.GetContents(ctx, input.Owner, input.Repo, path.Join(input.Path, file), &github.RepositoryContentGetOptions{
			Ref: input.Ref,
		})
		if err != nil {
			slogerr.WithError(v.logger, err).Debug("get action file")
			continue
		}
		if content == nil {
			// The path is a directory
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

func (v *validateStep) validate(ctx context.Context) error { //nolint:cyclop
	// Validate inputs
	// Docker container actions have no action file and no inputs to validate.
	if v.step.Uses == "" || strings.HasPrefix(v.step.Uses, "docker://") {
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
		return urfave.ErrSilent
	}
	return nil
}

func (v *validateStep) readLocalAction(action *workflow.Action) error {
	found := false
	for _, file := range []string{"action.yaml", "action.yml"} {
		p := filepath.Join(v.step.Uses, file)
		if err := workflow.ReadAction(p, action); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read a local action file: %w", err)
		}
		found = true
	}
	if !found {
		return errors.New("local action file not found")
	}
	return nil
}
