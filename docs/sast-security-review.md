# SAST Pipeline GitHub Actions security review

## Scope

This review covers `.github/workflows/sast.yml` in `ministryofjustice/github-actions` and its use by application repositories. The workflow performs dependency, code, infrastructure-as-code and publicly pullable container-image scans.

## Trust model

| Invocation | Secret-bearing scans | Behaviour |
|---|---:|---|
| Protected `main` execution | Allowed | Runs selected scans using credentials from an approved `main`-restricted central environment. |
| Other internal branch | Denied | Environment deployment policy rejects the job before scanner credentials are released. |
| Same-repository pull request | Denied | A no-secret boundary workflow records that credentialed scanning is deferred to protected execution. |
| Dependabot pull request | Denied | No repository-level scanner credentials exist and the protected environment rejects the ref. |
| Forked pull request | Denied | No repository-level scanner credentials exist and the protected environment rejects the ref. |
| `pull_request_target` | Denied | Validation fails explicitly. This prevents privileged credentials being combined with pull-request execution. |

Consumers name an approved protected central environment and do not pass scanner credentials. `secrets: inherit` is not part of the approved contract. Scanner credentials must not also exist at consumer repository level.

## Token permissions

The reusable workflow starts with `permissions: {}`. Permissions are granted per job:

| Job | Permissions | Reason |
|---|---|---|
| Validate inputs | `contents: read` | Check out and inspect caller paths and manifests. |
| Dependency scan | `contents: read` | Check out and scan caller content. |
| Code scan | `contents: read` | Check out and scan caller content. |
| IaC scan | `contents: read` | Check out and scan caller content. |
| Container scan | `contents: read` | Check out the optional Dockerfile and scan a publicly pullable image. |
| Publish results | `actions: read`, `contents: read`, `security-events: write` | Download validated artifacts and publish SARIF. |
| Policy result | none | Evaluate job outputs only. |
| Job summary | none | Write the workflow summary only. |

The primary workflow does not request `contents: write`, `id-token: write` or workflow-wide write permission. Private-registry authentication is outside this workflow and requires a separately reviewed workflow with its own permission contract.

## Third-party and external dependencies

| Dependency | Publisher | Purpose | Control |
|---|---|---|---|
| `actions/checkout` | GitHub | Check out caller content | Full commit SHA with release comment; credentials are not persisted. |
| `actions/upload-artifact` | GitHub | Retain validated scan results | Full commit SHA with release comment. |
| `actions/download-artifact` | GitHub | Retrieve validated SARIF | Full commit SHA with release comment. |
| `github/codeql-action/upload-sarif` | GitHub | Publish GitHub Code Scanning results | Full commit SHA with release comment. |
| Snyk CLI | Snyk | Execute the selected scans | Exact npm version; package integrity recorded during review. |

Dependabot monitors the `github-actions` ecosystem daily. Action update pull requests require normal review before merge. The Snyk CLI version is embedded in a shell installation command and therefore requires an explicit periodic version-and-integrity review until an automated approved update mechanism is introduced.

## Runner security

All jobs use GitHub-hosted `ubuntu-latest` runners. Each job receives a fresh hosted runner and jobs do not share a workspace or Docker daemon. The workflow does not support self-hosted runners.

If self-hosted support is proposed later, it requires a separate review of repository allow-lists, ephemeral execution, workspace cleanup, outbound network access, secret isolation and prohibition of untrusted pull requests.

## Artifact handling

| Property | Decision |
|---|---|
| Contents | Per-scan SARIF and minimal result metadata only. |
| Names | `sast-pipeline-<scan>-<run-id>-<attempt>`. |
| Retention | Seven days. |
| Access | Users and automation authorised to access the repository's Actions artifacts; Code Scanning access follows repository security permissions. |
| Classification | Treat as security-sensitive repository information; do not redistribute outside authorised repository and AppSec access. |
| Deletion | Automatic expiry after seven days; repository administrators may delete earlier for incident response or accidental sensitive content. |

The workflow does not upload source archives, complete workspaces, environment dumps or credentials. Only non-empty, structurally valid SARIF is eligible for publication.

## Residual risks and follow-up

- SARIF can contain file paths, finding descriptions and code-location context; access must remain limited to authorised repository and security users.
- A compromised scanner or action dependency remains a software-supply-chain risk; immutable SHA pins, release review and Dependabot reduce this risk.
- The exact Snyk CLI version requires a documented update cadence because Dependabot cannot update a version embedded in shell code.
- Credentialed pull-request scanning is deferred until a separately approved protected-environment or GitHub App design can scan untrusted changes without exposing long-lived credentials.
- The PoC environment currently allows administrator bypass. Production adoption must decide whether to disable bypass and add approved reviewers without creating an unavailable single-maintainer dependency.
- Private container registries require a separately reviewed workflow. They must not cause the primary workflow to request OIDC permissions for consumers that scan public images.

## Evidence

- PoC run `33312597313`: all four scans passed without `id-token: write`.
- PoC run `33312933607`: all four scans and SARIF publication passed using reviewed Node.js 24 artifact actions.
- PoC run `33330195801`: controlled presence-only test found repository-level credentials available to same-repository PR code; no values were disclosed.
- PoC run `33330786400`: protected environment rejected the feature ref before validation or scanning.
- PoC run `33330915899`: repeated presence-only test passed after repository secrets were removed and environment protection was applied.
- External actions are pinned to complete commit SHAs with readable release comments.
- Snyk CLI `1.1306.4` version, tarball and integrity value were recorded during review.

## Review conclusion

The primary workflow applies job-level least privilege, requires scanner credentials from an approved protected central environment, prevents protected credentials reaching pull-request-controlled code, rejects `pull_request_target`, uses GitHub-hosted runners, limits artifact scope and retention, and controls external dependencies with immutable pins and update review.
