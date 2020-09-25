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
      - uses: actions/checkout@v2
      - uses: ministryofjustice/github-actions/code-formatter@main
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

## Extending to other languages

If you want to add a code formatter for another language, you need to:

* Modify the `Dockerfile` to install the code formatting tool
* Modify `shared/code_formatter.rb` to
  * Identify files in the PR in the relevant language (by filename suffix)
  * Add a method to run the formatter, targeting each file
