package main

test_deny_job_should_have_permissions {
	not any_deny_job_should_have_permissions
}

any_deny_job_should_have_permissions {
	seeds := [
		{
			"exp": set(), "msg": "pass",
			"resource": {
				"path": ".github/workflows/release.yaml",
				"contents": yaml.unmarshal(`
jobs:
  release:
    permissions:
      contents: write
`),
			},
		},
		{
			"exp": {".github/workflows/test.yaml release: [GitHub Actions jobs should have permissions](https://github.com/suzuki-shunsuke/terraform-monorepo-github-actions/tree/main/policy/terraform/cloudwatch_log_retention_in_days.md)"},
			"msg": "no permissions",
			"resource": {
				"path": ".github/workflows/test.yaml",
				"contents": yaml.unmarshal(`
jobs:
  release:
    uses: suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@8e0d6d2a7171206b9d95b3b59fe74f8333b1be1b # v0.1.0
`),
			},
		},
	]

	some i
	seed := seeds[i]

	result := deny_job_should_have_permissions with input as seed.resource

	result != seed.exp
	trace(sprintf("FAIL %s (%d): %s, wanted %v, got %v", ["", i, seed.msg, seed.exp, result]))
}
