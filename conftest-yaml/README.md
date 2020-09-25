# Test all YAML with Conftest

A github action to apply [Conftest] policies written in Rego to all the YAML
files in a PR.

## USAGE

Create a file in your repo called `.github/workflows/conftest-yaml.yml` with the
following contents:

```
on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

name: Conftest YAML files
jobs:
  conftest-yaml:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: ministryofjustice/github-actions/conftest-yaml@main
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

By default, policies are assumed to be in the `./policy` directory of the
current repository (this is also the default behaviour of `conftest`).

You can override the policy directory by supplying a `POLICY_DIR` environment
variable.

You can pass additional command-line options to `conftest` in the
`CONFTEST_OPTIONS` environment variable.

[Conftest]: https://www.conftest.dev/
