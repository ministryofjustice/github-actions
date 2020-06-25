# YAML with escalated privileges

A Github Action to reject a PR if it contains values which will grant users escalated privileges.

Some values of YAML would grant users escalated privileges and if misapplied can affect other environment. This action is designed to detech such values before the corresponding PR is merged.

## Example

This is an example of such YAML, which will raise an error if
you try to parse it (at least, in ruby):

```
roleRef:
  kind: ClusterRole
  name: cluster-admin
```

## USAGE

Create a file in your repo called `.github/workflows/reject-escalated-privileges-yaml.yml` with the
following contents:

```
on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

jobs:
  reject-escalated-privileges-yaml:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - uses: ministryofjustice/github-actions/reject-escalated-privileges-yaml@main
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

