package policy_test

import (
	"log/slog"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/policy"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

func TestJobTimeoutMinutesIsRequiredPolicy_ApplyJob(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		job   *workflow.Job
		isErr bool
	}{
		{
			name: "normal",
			job: &workflow.Job{
				TimeoutMinutes: 30,
				Steps: []*workflow.Step{
					{
						Run: "echo hello",
					},
				},
			},
		},
		{
			name: "expression is used",
			job: &workflow.Job{
				TimeoutMinutes: "${{ matrix.timeout-minutes }}",
				Steps: []*workflow.Step{
					{
						Run: "echo hello",
					},
				},
			},
		},
		{
			name: "workflow using reusable workflow",
			job: &workflow.Job{
				Uses: "suzuki-shunsuke/renovate-config-validator-workflow/.github/workflows/validate.yaml@v0.2.3",
			},
		},
		{
			name: "job should have timeout-minutes",
			job: &workflow.Job{
				Steps: []*workflow.Step{
					{
						Run: "echo hello",
					},
				},
			},
			isErr: true,
		},
		{
			name: "all steps have timeout-minutes",
			job: &workflow.Job{
				Steps: []*workflow.Step{
					{
						Run:            "echo hello",
						TimeoutMinutes: 60,
					},
					{
						Run:            "echo hello",
						TimeoutMinutes: 60,
					},
				},
			},
		},
		{
			name: "expression is used in step's timeout-minutes",
			job: &workflow.Job{
				Steps: []*workflow.Step{
					{
						Run:            "echo hello",
						TimeoutMinutes: "${{ matrix.timeout-minutes }}",
					},
					{
						Run:            "echo hello",
						TimeoutMinutes: 60,
					},
				},
			},
		},
	}
	p := &policy.JobTimeoutMinutesIsRequiredPolicy{}
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if err := p.ApplyJob(logger, nil, nil, d.job); err != nil {
				if !d.isErr {
					t.Fatal(err)
				}
				return
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
