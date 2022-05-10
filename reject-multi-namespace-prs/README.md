# Reject Multi-namespace PRs

Pull requests (PRs) against the [environments repository][env-repo],
raised by users of the [MoJ Cloud Platform][cloud-platform],
have to be approved by the Cloud Platform team.

If these PRs affect multiple namespaces, they can be quite large,
consisting of kubernetes YAML files and terraform code. This can
be hard to review correctly, making it easier for mistakes to slip
through.

This Github Action marks PRs as failed if they affect more than
one namespace folder.

## USAGE

Create a file in your repo called `.github/workflows/reject-multi-namespace-prs.yml` with the
following contents:

```
on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

jobs:
  reject-multi-namespace-prs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: ministryofjustice/github-actions/reject-multi-namespace-prs@main
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

[env-repo]: https://github.com/ministryofjustice/cloud-platform-environments
[cloud-platform]: https://github.com/ministryofjustice/cloud-platform
