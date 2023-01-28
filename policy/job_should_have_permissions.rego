package main

deny_job_should_have_permissions[msg] {
	# jobs.<job name>.permissions
	some jobName
	job := input.contents.jobs[jobName]
	not job.permissions
	msg = sprintf("%s %s: [GitHub Actions jobs should have permissions](%s)", [input.path, jobName, "https://github.com/suzuki-shunsuke/terraform-monorepo-github-actions/tree/main/policy/terraform/cloudwatch_log_retention_in_days.md"])
}
