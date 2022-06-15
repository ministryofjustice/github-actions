# PR Code Formatter

A Github Action to apply code formatting to file in **PRs only**. CI GH Action will run and automatically commit to the same branch when there is a difference in files in the PR.

Formats Ruby, Terraform, YAML/YML, Python, Markdown and JSON file within a PR.

## USAGE

Create a file in your repo called `.github/workflows/format-code.yml` with the
following contents:

```
on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

jobs:
  format-code:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: ministryofjustice/github-actions/code-formatter@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.
