# Terraform Static Analysis Action

This action combines [TFSEC](https://github.com/tfsec/tfsec), [Checkov](https://github.com/bridgecrewio/checkov) and [tflint](https://github.com/terraform-linters/tflint) into one action, loosely based on the [TFSEC action](https://github.com/triat/terraform-security-scan) and [Checkov actions](https://github.com/bridgecrewio/checkov-action) here.

The main reason for combining these checks is to enable one action to run which can cover multiple checks as well as multiple and nested Terraform folders.  This action also has logic to perform different scan options depending if you want to scan your whole repo or only individual or changed folders:

Full scan (`full`) - scan all folders with `*.tf` files in a repository.

Changes only (`changed`) - scan only folders with `*.tf` files that have had changes since the last commit.

Single folder (`single`) - standard scan of a given folder.

See the [action.yml](action.yml) for other input options. Global excludes for checks can be added at the action level, or inline exclude comments can be added for each check (see the check's user guide for correct syntax).

## Example

```
jobs:
  terraform-static-analysis:
    name: Terraform Static Analysis
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v2.3.4
      with:
        fetch-depth: 0
    - name: Run Analysis
      uses: ministryofjustice/github-actions/terraform-static-analysis@main
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        AWS_ACCESS_KEY_ID:  ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY:  ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        TF_IN_AUTOMATION: true
      with:
        scan_type: changed
        tfsec_exclude: AWS095
```

### Notes

`fetch-depth: 0` is required to get a git diff to detect changed files.

`GITHUB_TOKEN` is required to write the results to the pull request.

`AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and TF_IN_AUTOMATION` are required to run Terraform init, this enables TFSec to run tests on called modules as well.
