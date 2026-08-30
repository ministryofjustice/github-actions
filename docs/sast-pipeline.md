# SAST Pipeline consumer contract

The SAST Pipeline is a reusable GitHub Actions workflow maintained by OCTO/AppSec. It currently uses Snyk internally, but its consumer interface is product-neutral.

Location: `.github/workflows/sast.yml`

## Minimal caller

Reference a reviewed release tag or full commit SHA. Do not reference `main` or another mutable branch in production.

```yaml
name: Security scanning

on:
  workflow_dispatch:
  schedule:
    - cron: "0 5 * * *"

permissions: {}

jobs:
  sast:
    permissions:
      actions: read
      contents: read
      security-events: write
    uses: ministryofjustice/github-actions/.github/workflows/sast.yml@<full-commit-sha>
    with:
      run_dependency_scan: true
      run_code_scan: true
      run_iac_scan: true
      run_container_scan: false
      working_directory: "."
      iac_path: "infrastructure"
      severity_threshold: high
      upload_sarif: true
      fail_build: false
    secrets:
      scanner_client_id: ${{ secrets.SNYK_CLIENT_ID }}
      scanner_client_secret: ${{ secrets.SNYK_CLIENT_SECRET }}
```

The called workflow cannot elevate permissions supplied by the caller.

## Inputs

| Input | Type | Default | Meaning and validation |
|---|---|---:|---|
| `run_dependency_scan` | boolean | `true` | Scan supported dependency manifests under `working_directory`. Fails when no supported manifest exists. A `package.json` at `working_directory` requires a supported npm, Yarn or pnpm lock file. |
| `run_code_scan` | boolean | `true` | Scan source under `working_directory`. Requires Snyk Code entitlement. |
| `run_iac_scan` | boolean | `false` | Scan only `iac_path`. Fails when no supported IaC file exists. |
| `run_container_scan` | boolean | `false` | Scan an existing image. The workflow never builds or pushes it. |
| `working_directory` | string | `.` | Existing repository-relative directory. Absolute paths and traversal are rejected. |
| `iac_path` | string | `.` | Existing directory relative to `working_directory`. |
| `dockerfile_path` | string | empty | Optional repository-relative Dockerfile providing container context. |
| `image_uri` | string | empty | Required for a container scan. The image must be publicly pullable and an immutable digest is strongly preferred. |
| `severity_threshold` | string | `high` | One of `low`, `medium`, `high` or `critical`. Code maps `critical` to `high` because Snyk Code has no critical severity. |
| `upload_sarif` | boolean | `true` | Upload only validated SARIF to GitHub Code Scanning. |
| `fail_build` | boolean | `false` | Fail when findings meet the threshold. Operational errors always fail. |

At least one scan must be selected. Container-only inputs are rejected unless the container scan is enabled.

## Secrets

Pass secrets explicitly. Do not use `secrets: inherit`.

| Secret | Required | Purpose |
|---|---:|---|
| `scanner_client_id` | conditional | Client ID for the approved non-personal OAuth service account. Must be paired with `scanner_client_secret`. |
| `scanner_client_secret` | conditional | Client secret for the approved non-personal OAuth service account. |
| `scanner_token` | conditional | Non-personal static-token fallback. Do not provide it together with the OAuth pair. |

Trusted runs require either the OAuth pair or `scanner_token`. Secret-bearing scans do not run for forked or Dependabot pull requests.

## Permissions

| Permission | When required |
|---|---|
| `contents: read` | Checkout and all repository scans. |
| `actions: read` | Download validated artifacts before SARIF publication. |
| `security-events: write` | Upload SARIF to GitHub Code Scanning. |

The workflow does not require `contents: write` or `id-token: write`. Private-registry authentication requires a separately reviewed workflow and permission contract.

## Results and failure rules

- A clean result means the selected scan executed without findings at the threshold.
- Findings are recorded separately from scanner errors.
- Unsupported or missing projects are operational failures, not clean scans.
- Findings fail the workflow only when `fail_build` is `true`.
- Operational, authentication, invalid-result and requested-publication failures always fail.
- Only existing, non-empty, structurally valid SARIF is uploaded, using stable per-scan categories.
- Result artifacts are retained for seven days for diagnosis.
- The final summary records requested, executed and skipped scans and the policy result.

Version 1 intentionally exposes no outputs and never outputs credentials or raw scanner evidence.
