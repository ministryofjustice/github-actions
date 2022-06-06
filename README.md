# Github Actions

[![repo standards badge](https://img.shields.io/badge/dynamic/json?color=blue&style=for-the-badge&logo=github&label=MoJ%20Compliant&query=%24.data%5B%3F%28%40.name%20%3D%3D%20%22github-actions%22%29%5D.status&url=https%3A%2F%2Foperations-engineering-reports.cloud-platform.service.justice.gov.uk%2Fgithub_repositories)](https://operations-engineering-reports.cloud-platform.service.justice.gov.uk/github_repositories#github-actions "Link to report")

A collection of github actions.

<!-- markdownlint-disable MD013 -->
| Action       | Description                                              |
|--------------|----------------------------------------------------------|
| [branch-deleter](branch-deleter) | Delete a branch after its PR is merged |
| [code-formatter](code-formatter) | Run various code formatters against a PR, and commit the results |
| [conftest-yaml](conftest-yaml) | Use [Conftest] to check your YAML against Rego policies |
| [iam-role-policy-changes-check](iam-role-policy-changes-check) |  |
| [malformed-yaml](malformed-yaml) | Reject a PR if it contains malformed YAML files |
| [reject-multi-namespace-prs](reject-multi-namespace-prs) | Reject a PR if it affects multiple namespace folders |
| [terraform-static-analysis](terraform-static-analysis) | Runs TFSec, Checkov and TFlint against Terraform |
<!-- markdownlint-enable MD013 -->

[Conftest]: https://www.conftest.dev/
