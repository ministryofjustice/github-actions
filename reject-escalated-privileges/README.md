# Malformed YAML

A Github Action to reject a PR if it contains malformed YAML files.

YAML syntax errors can be hard to spot by eye, and a malformed file
can break a CI/CD pipeline. This action is designed to prevent such
errors before the corresponding PR is merged.

## Example

This is an example of malformed YAML, which will raise an error if
you try to parse it (at least, in ruby):

```
desc: Example of a malformed YAML file
key: "\"
```

## USAGE

Create a file in your repo called `.github/workflows/malformed-yaml.yml` with the
following contents:

```
on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

jobs:
  reject-malformed-yaml:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - uses: ministryofjustice/github-actions/malformed-yaml@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

