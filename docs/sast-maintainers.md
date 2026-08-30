# SAST Pipeline maintainer guide

## Ownership and architecture

The reusable workflow is `.github/workflows/sast.yml` in `ministryofjustice/github-actions`. Consumers call it by immutable SHA or reviewed release tag. The public interface is product-neutral; Snyk is the initial implementation.

Execution flow:

1. `validate-inputs` rejects unsafe contexts, paths and unsupported combinations.
2. `dependency-scan`, `code-scan`, `iac-scan` and `container-scan` run independently when selected.
3. Each scan classifies clean, findings or operational error and validates its SARIF.
4. `publish-results` downloads validated artifacts and publishes stable categories.
5. `policy-result` combines validation, scan and publication states.
6. `job-summary` records the final human-readable outcome.

Jobs use fresh GitHub-hosted runners, explicit timeouts and job-level least privilege. The primary container scan accepts public images only and does not build, push or authenticate to a registry.

## Snyk organisation and credentials

The current workflow targets Snyk organisation `octo-Sq42AMT8YK8yMg4AtYeKJ5`. Changing it is a security-impacting configuration change and requires Snyk-administrator and OCTO/AppSec review.

Use a non-personal OAuth client with least privilege. A static non-personal token is a fallback only when explicitly approved. Credentials belong to the approved central protected environment, not consumer repositories. The environment owner controls creation, rotation, revocation, trusted refs and reviewer policy.

Central protected-environment provisioning and cross-repository behavior must be proven before release. Never work around a missing environment by creating repository-level scanner secrets or using `secrets: inherit`.

## Release process

1. Work on a short-lived branch and keep changes scoped to the pipeline.
2. Update workflow, consumer contract, security review, policy and ADRs together when applicable.
3. Run `actionlint`, `git diff --check` and the PoC functional/negative suites.
4. Record immutable-SHA PoC runs, Code Scanning evidence and rollback evidence.
5. Obtain the reviews listed in the release PR.
6. Merge without bypassing repository controls.
7. Create an annotated release tag on the reviewed merge commit using the repository's agreed semantic version.
8. Verify the tag resolves to the intended commit and rerun the PoC using the tag.
9. Publish release notes describing contract changes, security impact, migration and rollback.

No release tag is approved while Stage 9 production prerequisites or required reviews remain outstanding.

## Tag management

- Treat release tags as immutable; never move or force-update a published tag.
- Prefer a full SHA for pre-release testing and high-assurance consumers.
- Use semantic versioning: patch for compatible fixes, minor for compatible inputs/features and major for breaking contract changes.
- If a tag is wrong, publish a new corrected version and deprecate the bad version.
- Verify locally and remotely with `git rev-parse <tag>^{}` and `git ls-remote --tags origin`.

## Testing

The PoC repository is `ministryofjustice/sast-pipeline-poc`. Required evidence includes individual scans, all scans, disabled options, custom paths/image, clean fixtures, negative cases, PR boundaries, immutable SHA/tag and rollback. Follow `docs/sast-poc-test-report.md` in the PoC repository. A skipped or unavailable scan is not a passing clean test.

## Dependency updates

- GitHub Actions remain pinned to full SHAs with readable version comments.
- Review upstream release notes and publisher ownership before updating.
- Snyk CLI remains pinned to an exact npm version; verify version, tarball URL and integrity.
- Run all Stage 9 regression scenarios after an action, CLI or SARIF uploader update.
- Do not merge automated update PRs without normal review and passing evidence.

## Rollback

Keep the previous known-good immutable workflow SHA in release notes. To roll back, update the consumer `uses` reference to that SHA, validate YAML, run the supported scans and record the successful run. Rollback must not restore repository secrets, PR credential access, OIDC or private-registry permissions.

## Incident response

For suspected credential exposure, unsafe PR execution, compromised dependency or incorrect clean result:

1. Pause rollout and enforcement; do not delete evidence.
2. Revoke/rotate affected scanner credentials through the environment owner.
3. Disable or restrict the affected environment if needed.
4. Preserve run URLs, job IDs, workflow SHA and artifact metadata without copying secrets.
5. Notify the central pipeline owner, Snyk administrator and OCTO/AppSec.
6. Roll consumers back to the last known-good SHA if safe.
7. Fix, independently review and rerun Stage 9 before resuming.
8. Record cause, impact, remediation and follow-up.

## Deprecation

- Announce a replacement, reason, affected versions, migration guide, deadline and support contact.
- Allow a documented migration window agreed with service owners.
- Do not delete a release tag or silently change its behavior.
- Security-critical versions may be blocked sooner only with an incident decision and explicit communication.
- Retain rollback documentation until all known consumers have migrated.

## Support and review ownership

`@ministryofjustice/github-community` owns repository workflow review. OCTO/AppSec or its delegated owner reviews security policy and protected-environment changes. The LAA maintainer reviews migration accuracy, and the Snyk administrator reviews organisation, entitlement and credentials.
