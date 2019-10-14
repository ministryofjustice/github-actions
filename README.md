# Github Action: Code Formatter

A Github Action to apply code formatters to PRs, and commit any resulting changes.

The following formatters will be applied:

* `*.tf` files -> `terraform fmt`
* `*.rb` files -> `standardrb --fix`

## USAGE

```
uses: ministryofjustice/github-action-code-formatter@v1
```

