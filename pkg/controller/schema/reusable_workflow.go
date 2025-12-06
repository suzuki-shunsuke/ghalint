package schema

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"gopkg.in/yaml.v3"
)

func (v *validateJob) validateReusableWorkflow(ctx context.Context) error {
	// read workflow
	wf := &ReusableWorkflow{}
	if err := v.read(ctx, wf); err != nil {
		return fmt.Errorf("read a reusable workflow: %w", err)
	}
	if err := v.validateWorkflow(wf); err != nil {
		return fmt.Errorf("validate a reusable workflow: %w", err)
	}
	return nil
}

/*
on:
  workflow_call:
    inputs:
      aqua_policy_config:
        required: false
        type: string
*/

type ReusableWorkflow struct {
	On *On
}

type On struct {
	WorkflowCall *WorkflowCall `yaml:"workflow_call"`
}

func (o *On) UnmarshalYAML(unmarshal func(any) error) error { //nolint:cyclop
	var onAny any
	if err := unmarshal(&onAny); err != nil {
		return fmt.Errorf("unmarshal a workflow to any: %w", err)
	}
	if s, ok := onAny.(string); ok {
		if s != "workflow_call" {
			return nil
		}
		o.WorkflowCall = &WorkflowCall{}
		return nil
	}
	onMap, ok := onAny.(map[string]any)
	if !ok {
		return errors.New("failed to convert workflow on into map")
	}
	workflowCallAny, ok := onMap["workflow_call"]
	if !ok {
		return nil
	}
	o.WorkflowCall = &WorkflowCall{}
	workflowCallMap, ok := workflowCallAny.(map[string]any)
	if !ok {
		return nil
	}
	inputsAny, ok := workflowCallMap["inputs"]
	if !ok {
		return nil
	}
	inputsMap, ok := inputsAny.(map[string]any)
	if !ok {
		return nil
	}
	o.WorkflowCall.Inputs = map[string]*workflow.Input{}
	for inputKey, v := range inputsMap {
		o.WorkflowCall.Inputs[inputKey] = &workflow.Input{}
		inputValueMap, ok := v.(map[string]any)
		if !ok {
			continue
		}
		requiredAny, ok := inputValueMap["required"]
		if !ok {
			continue
		}
		required, ok := requiredAny.(bool)
		if !ok {
			continue
		}
		o.WorkflowCall.Inputs[inputKey] = &workflow.Input{
			Required: required,
		}
	}
	return nil
}

type WorkflowCall struct {
	Inputs map[string]*workflow.Input
}

func (v *validateJob) validateWorkflow(wf *ReusableWorkflow) error {
	if wf.On == nil {
		return errors.New("the reusable workflow is invalid. on is not set")
	}
	if wf.On.WorkflowCall == nil {
		return errors.New("the reusable workflow is invalid. workflow_call is not set")
	}
	inputs := wf.On.WorkflowCall.Inputs
	requiredKeys := map[string]struct{}{}
	for key, input := range inputs {
		if input.Required {
			requiredKeys[key] = struct{}{}
		}
	}
	v.logger = v.logger.With(
		"valid_inputs", strings.Join(slices.Collect(maps.Keys(inputs)), ", "),
		"required_inputs", strings.Join(slices.Collect(maps.Keys(requiredKeys)), ", "),
	)
	failed := false
	// Check if the input is valid
	for key := range v.job.With {
		if _, ok := inputs[key]; !ok {
			v.logger.Error("invalid input key", "input_key", key)
			failed = true
		}
	}
	// Check if required keys are set
	for key := range requiredKeys {
		if _, ok := v.job.With[key]; !ok {
			v.logger.Error("required key is not set", "input_key", key)
			failed = true
		}
	}
	if failed {
		return urfave.ErrSilent
	}
	return nil
}

func readReusableWorkflow(fs afero.Fs, p string, wf *ReusableWorkflow) error {
	f, err := fs.Open(p)
	if err != nil {
		return fmt.Errorf("open a workflow file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(wf); err != nil {
		err := fmt.Errorf("parse a workflow file as YAML: %w", err)
		if errors.Is(err, io.EOF) {
			return slogerr.With(err, //nolint:wrapcheck
				"reference", "https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/codes/001.md",
			)
		}
		return err
	}
	return nil
}

func (v *validateJob) read(ctx context.Context, wf *ReusableWorkflow) error { //nolint:cyclop
	if strings.HasPrefix(v.job.Uses, "./") {
		// local workflow
		if err := readReusableWorkflow(v.fs, v.job.Uses, wf); err != nil {
			return fmt.Errorf("read a local workflow file: %w", err)
		}
		return nil
	}
	// <owner>/<repo>[/<path>]@<ref>
	fullPath, ref, ok := strings.Cut(v.job.Uses, "@")
	if !ok {
		return fmt.Errorf("invalid job.uses format: %s", v.job.Uses)
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
	cachePath := filepath.Join(v.rootDir, "actions", owner, repo, sha, path)
	if f, err := afero.Exists(v.fs, cachePath); err != nil {
		return fmt.Errorf("check if the workflow file exists: %w", err)
	} else if f {
		if err := readReusableWorkflow(v.fs, cachePath, wf); err != nil {
			return fmt.Errorf("read a cached workflow file: %w", err)
		}
		return nil
	}
	// Download a wofklow file
	content, _, _, err := v.gh.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: sha,
	})
	if err != nil {
		return fmt.Errorf("download workflow file: %w", err)
	}
	// write workflow to the cache dir
	if err := v.fs.MkdirAll(filepath.Dir(cachePath), dirPermission); err != nil {
		return fmt.Errorf("create workflow directory: %w", err)
	}
	c, err := content.GetContent()
	if err != nil {
		return fmt.Errorf("get content: %w", err)
	}
	b := []byte(c)
	if err := afero.WriteFile(v.fs, cachePath, b, filePermission); err != nil {
		return fmt.Errorf("write workflow file: %w", err)
	}
	if err := yaml.Unmarshal(b, wf); err != nil {
		return fmt.Errorf("unmarshal workflow file: %w", err)
	}
	return nil
}
