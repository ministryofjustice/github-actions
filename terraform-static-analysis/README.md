# Terraform Static Analysis Action

This action combines [TFSEC](https://github.com/tfsec/tfsec), [Checkov](https://github.com/bridgecrewio/checkov) and [tflint](https://github.com/terraform-linters/tflint) into one action, loosely based on the [TFSEC action](https://github.com/triat/terraform-security-scan) and [Checkov actions](https://github.com/bridgecrewio/checkov-action) here.

The main reason for combining these checks is to enable one action to run which can cover multiple checks as well as multiple and nested Terraform folders.  This action also has logic to perform different scan options depending if you want to scan your whole repo or only individual or changed folders:

Full scan (`full`) - scan all folders with `*.tf` files in a repository.

Changes only (`changed`) - scan only folders with `*.tf` files that have had changes since the last commit.

Single folder (`single`) - standard scan of a given folder.

See the [action.yml](action.yml) for other input options. Global excludes for checks can be added at the action level, or inline exclude comments can be added for each check (see the check's user guide for correct syntax). Additional tflint configurations can also be added to the `tflint-configs` directory and then passed with the appropriate input option.

With the introduction of trivy as a scanner option you can now select between using tfsec or trivy by adding the input tfsec_trivy: leaving this blank will default the scanner to tfsec adding trivy to it will switch the scanner to trivy. If you decide to use trivy you will need to add the following inputs to the with statement

If Trivy is enabled and you want to ignore some errors a trivy ignore file will need to be created in your repo and the path to that file added to the input of trivy_ignore:, an example of this file is shown bellow.


## Example trivy config
```
with:
        scan_type: full
        tfsec_trivy: trivy
        trivy_severity: HIGH,CRITICAL
        trivy_ignore: ./.trivyingnore.yaml
        checkov_exclude: CKV_GIT_1,CKV_AWS_126,CKV2_AWS_38,CKV2_AWS_39
        tflint_exclude: terraform_unused_declarations
        tflint_call_module_type: none


```

## Example

```
jobs:
  terraform-static-analysis:
    name: Terraform Static Analysis
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v3
      with:
        fetch-depth: 0
    - name: Run Analysis
      uses: ministryofjustice/github-actions/terraform-static-analysis@main
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with:
        scan_type: changed
        tfsec_exclude: AWS095
```

## Example trivy ignore file

```
vulnerabilities:
  - id: CVE-2022-40897
    paths:
      - "usr/local/lib/python3.9/site-packages/setuptools-58.1.0.dist-info/METADATA"
    statement: Accept the risk
  - id: CVE-2023-2650
  - id: CVE-2023-3446
  - id: CVE-2023-3817
  - id: CVE-2023-29491
    expired_at: 2023-09-01

misconfigurations:
  - id: AVD-DS-0001
  - id: AVD-DS-0002
    paths:
      - "docs/Dockerfile"
    statement: The image needs root privileges

secrets:
  - id: aws-access-key-id
  - id: aws-secret-access-key
    paths:
      - "foo/bar/aws.secret"

licenses:
  - id: GPL-3.0 # License name is used as ID
    paths:
      - "usr/share/gcc/python/libstdcxx/v6/__init__.py"

```


### Notes

`fetch-depth: 0` is required to get a git diff to detect changed files.

`GITHUB_TOKEN` is required to write the results to the pull request. This is the built in workflow token created when you start using Actions (see [here](https://docs.github.com/en/actions/reference/authentication-in-a-workflow)) this should have read and write permissions to write to a pull request, this can be found under `Actions` in the repository `Settings`
