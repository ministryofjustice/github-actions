# Github Actions

[![repo standards badge](https://img.shields.io/badge/dynamic/json?color=blue&flat-square&logo=github&label=MoJ%20Compliant&query=%24.result&url=https%3A%2F%2Foperations-engineering-reports.cloud-platform.service.justice.gov.uk%2Fapi%2Fv1%2Fcompliant_public_repositories%2Fgithub-actions)](https://operations-engineering-reports.cloud-platform.service.justice.gov.uk/public-github-repositories.html#github-actions "Link to report")

A collection of github actions.

<!-- markdownlint-disable MD013 -->
| Action       | Description                                              |
|--------------|----------------------------------------------------------|
| [code-formatter](code-formatter) | Run various code formatters against a PR, and commit the results |
| [iam-role-policy-changes-check](iam-role-policy-changes-check) | Reject PR if it contains IAM related content |
| [malformed-yaml](malformed-yaml) | Reject a PR if it contains malformed YAML files |
| [reject-multi-namespace-prs](reject-multi-namespace-prs) | Reject a PR if it affects multiple namespace folders |
| [terraform-static-analysis](terraform-static-analysis) | Runs TFSec, Checkov and TFlint against Terraform |
<!-- markdownlint-enable MD013 -->

[Conftest]: https://www.conftest.dev/
