# SAST Pipeline consumer contract

Rollout, enforcement and exception rules are defined in the [SAST enforcement and exception policy](sast-enforcement-policy.md).

The SAST Pipeline is a reusable GitHub Actions workflow maintained by OCTO/AppSec. It currently uses Snyk internally, but its consumer interface is product-neutral.

Location: `.github/workflows/sast.yml`

## What it scans

- supported dependency manifests and lock files under `working_directory`;
- application source with Snyk Code;
- Terraform, YAML and JSON IaC under `iac_path`;
- an existing, publicly pullable container image supplied by immutable digest where possible.

## What it does not scan

- running applications or HTTP endpoints (DAST);
- secrets committed to source;
- licences or repository compliance policy;
- images that must be built or pushed by this workflow;
- private registries in the primary workflow;
- self-hosted runners;
- pull-request-controlled code with scanner credentials.

Use the appropriate separately approved control for these cases. Do not add broad permissions to this workflow to make an unsupported use case work.

## Prerequisites

- GitHub Actions and GitHub Code Scanning are enabled for the consumer repository.
- The caller grants the documented job permissions.
- Selected paths, manifests, lock files and image references exist.
- Snyk Code is enabled for the target Snyk organisation when code scanning is selected.
- The centrally approved protected scanner environment and credentials have been provisioned.
- The caller references an approved immutable SHA or release tag.

Protected-environment provisioning is a pre-release administrative dependency. Until the central environment design is approved and proven, consumers must not create repository-level scanner secrets as a workaround.

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
      scanner_environment: "sast-pipeline-production"
```

The called workflow cannot elevate permissions supplied by the caller.

Every job has an explicit timeout. Validation and publication are limited to 10 minutes, dependency and IaC scans to 15 minutes, code and container scans to 20 minutes, and policy/summary jobs to 5 minutes. A timed-out selected scan is an operational failure and must not be reported as clean.

The approved central environment must contain `SCANNER_CLIENT_ID` and `SCANNER_CLIENT_SECRET`, must be restricted to trusted refs, and must not expose those values as consumer repository-level secrets. The approved PoC design permits only `main`. Pull-request workflows do not call the credentialed scanner; they record that scanning is deferred to protected post-merge or authorised manual execution.

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
| `scanner_environment` | string | none | Required approved central environment containing scanner credentials. Restrict deployment to trusted refs such as `main`. |

At least one scan must be selected. Container-only inputs are rejected unless the container scan is enabled.

## Secrets and protected environment

Do not pass scanner credentials from the caller and do not use `secrets: inherit`. Store credentials only as secrets in the protected environment named by `scanner_environment`.

| Secret | Required | Purpose |
|---|---:|---|
| `SCANNER_CLIENT_ID` | conditional | Client ID for the approved non-personal OAuth service account. Must be paired with `SCANNER_CLIENT_SECRET`. |
| `SCANNER_CLIENT_SECRET` | conditional | Client secret for the approved non-personal OAuth service account. |
| `SCANNER_TOKEN` | conditional | Non-personal static-token fallback. Do not provide it together with the OAuth pair. |

Trusted runs require either the OAuth pair or `SCANNER_TOKEN` in the selected protected environment. Do not duplicate these values as consumer repository-level secrets. Credentialed scans must not run from pull-request-controlled workflow code, including same-repository, forked and Dependabot pull requests.

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

## Result locations

| Result | Location |
|---|---|
| Human summary | GitHub Actions run summary |
| Job state and reason | Individual scan and `security / Policy result` jobs |
| Validated SARIF | Actions artifact named `sast-pipeline-<scan>-<run-id>-<attempt>` for seven days |
| Published findings | Repository **Security → Code scanning** |
| Scanner-side project data | Approved Snyk organisation, subject to scanner administration and retention |

## Onboarding

1. Confirm the repository and selected project types are supported.
2. Obtain approval to use the central protected scanner environment.
3. Copy the minimal caller into `.github/workflows/security.yml`.
4. Pin `uses` to the approved full SHA or release tag.
5. Start with `severity_threshold: high` and `fail_build: false`.
6. Run an authorised manual or protected-branch test.
7. Confirm every selected scan executed, inspect the summary and Code Scanning results, and retain the run URL.
8. Add the stable policy check to branch protection only when the approved enforcement release requires it.

Migration from the LAA workflow is documented in [SAST Pipeline migration guide](sast-migration.md).

## Troubleshooting

| Symptom | Likely cause | Action |
|---|---|---|
| Workflow fails before jobs start | Caller permissions are lower than the reusable workflow requests | Grant only the documented `actions`, `contents` and `security-events` permissions |
| `input_validation_failed` | Invalid path, unsupported project, missing lock file or missing image | Read the Validate inputs annotation and correct the caller input/repository fixture |
| `scanner authentication is required` | Protected environment or credentials are unavailable to the trusted ref | Ask the central environment owner; do not add repository secrets |
| `scan_operational_error` | Scanner, entitlement, network, image pull or timeout failure | Inspect the selected scan job and rerun after correcting the operational cause |
| `result_publication_failed` | SARIF upload failed or permission was denied | Confirm `security-events: write`, Code Scanning availability and artifact validity |
| Scan is skipped | The scan toggle is false or the context is deliberately untrusted | Check the summary; skipped is not clean |
| Clean Code scan has no artifact | Snyk Code returned no findings and no SARIF file | Confirm the Code job executed and its status is clean |
| Private image cannot be pulled | Private registries are unsupported here | Use a separately approved private-registry workflow design |

## Support

- Implementation and workflow questions: open an issue or pull request in `ministryofjustice/github-actions` and notify `@ministryofjustice/github-community`.
- Scanner organisation, entitlement or credential questions: contact the approved Snyk administrator/OCTO AppSec owner.
- Suspected credential exposure or unsafe workflow behavior: stop the rollout, preserve the run URL without copying secrets, and follow the maintainer incident process.

Maintainers should use the [maintainer guide](sast-maintainers.md). Enforcement decisions are in the [enforcement policy](sast-enforcement-policy.md).
