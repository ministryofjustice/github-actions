# Migrate from the LAA SAST workflow

## Scope

This guide migrates callers from:

`ministryofjustice/laa-reusable-github-actions/.github/workflows/sast.yml@<sha>`

to:

`ministryofjustice/github-actions/.github/workflows/sast.yml@<approved-sha-or-tag>`

The central workflow is not a drop-in replacement. Review every mapping and unsupported behavior.

## Contract mapping

| LAA input/secret | Central equivalent | Migration action |
|---|---|---|
| `minimum_threshold` | `severity_threshold` | Rename; keep an approved value |
| implicit always-on Code | `run_code_scan: true` | Set explicitly |
| no dependency scan | `run_dependency_scan` | Select explicitly and provide a supported manifest/lock file |
| non-empty `iac_path` enables IaC | `run_iac_scan: true` plus `iac_path` | Set both |
| non-empty `image_uri` enables container | `run_container_scan: true` plus `image_uri` | Set both and use a public immutable digest |
| `dockerfile_path` default `Dockerfile` | `dockerfile_path` default empty | Set when container context is required |
| `exclude_base_image_vulns` | none | Not supported; central workflow reports full container evidence |
| `ecr_session_name` | none | Remove |
| `ECR_ROLE_TO_ASSUME` | none | Remove; private ECR is outside the primary workflow |
| caller `SNYK_CLIENT_ID`/`SNYK_CLIENT_SECRET` | centrally protected environment | Remove caller mapping; do not create repository secrets |
| caller `SNYK_TOKEN` | centrally protected environment fallback | Remove caller mapping |
| no equivalent | `working_directory` | Set for monorepos; otherwise use `.` |
| no equivalent | `upload_sarif` | Normally `true` |
| no equivalent | `fail_build` | Release 1 uses `false` |
| no equivalent | `scanner_environment` | Use only the centrally approved environment name |

## Permission changes

Remove LAA's `id-token: write`. The central primary workflow requires only:

```yaml
permissions:
  actions: read
  contents: read
  security-events: write
```

Do not add OIDC merely to preserve private-ECR behavior. Use a separately approved design if private-image scanning remains required.

## Example central caller

```yaml
name: Security scanning

on:
  workflow_dispatch:
  push:
    branches: [main]

permissions: {}

jobs:
  security:
    permissions:
      actions: read
      contents: read
      security-events: write
    uses: ministryofjustice/github-actions/.github/workflows/sast.yml@<approved-full-sha>
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
      scanner_environment: "<approved-central-environment>"
```

Credentialed scanning must not execute in PR-controlled workflow context. Use the approved no-secret PR boundary and defer scanning to protected execution.

## Validation

1. Validate workflow syntax with `actionlint` and `git diff --check`.
2. Confirm no caller Snyk/ECR secret mappings or `id-token: write` remain.
3. Confirm the `uses` reference is immutable.
4. Run each selected scan and then all scans together.
5. Confirm clean fixtures are clean rather than skipped.
6. Confirm controlled findings and Code Scanning categories.
7. Confirm operational errors fail honestly and logs contain no credentials.
8. Record run URLs and obtain the required review.

## Rollback

Before migration, record the exact LAA SHA currently used. If the central rollout must be reversed:

1. restore that reviewed LAA SHA and its former input/permission contract;
2. do not restore repository-level scanner secrets without an explicit security decision;
3. validate and run the supported scans;
4. record the rollback run and incident/change reference;
5. notify central and LAA workflow owners.

If restoring the old credential model would recreate a known exposure, pause delivery or use an approved emergency control instead of an unsafe rollback.

## Deprecation

The LAA deprecation deadline is **to be agreed** by the LAA maintainer, central workflow owner and affected services after the central production regression and release tag succeed. Until then, migrations are controlled opt-ins. The deadline, final supported LAA SHA and support window must be announced before enforcement.
