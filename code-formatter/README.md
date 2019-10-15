# Code Formatter

A Github Action to apply code formatters to PRs, and commit any resulting changes.

The following formatters will be applied:

* `*.tf` files -> `terraform fmt`
* `*.rb` files -> `standardrb --fix`

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
      - uses: actions/checkout@master
      - uses: ministryofjustice/github-actions/code-formatter@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

