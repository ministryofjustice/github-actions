# SAST enforcement and exception policy

## Status and purpose

| Field | Value |
|---|---|
| Owner | Central SAST pipeline owners |
| Initial rollout | Release 1 — Observe |
| Approval | Proposed for service-owner and OCTO/AppSec review |
| Review | At each promotion and at least every six months |

The scanner identifies issues; this policy decides whether findings or operational failures block delivery. It covers dependency, code, infrastructure-as-code (IaC) and container scanning.

The stable policy check is `security / Policy result`. Individual optional scan jobs must not be configured as separate required checks.

## Principles

- A scan that did not complete must never be reported as clean.
- Existing findings must not unexpectedly block every consumer during rollout.
- New high-risk findings are the intended enforcement target.
- Operational and publication failures remain distinct from findings.
- Exceptions must be explicit, approved, time-limited and auditable.
- Enforcement requires matching Stage 9 evidence.

## Scan decisions

| Scan | Release 1 | Release 2 | Release 3 threshold | Existing findings | New findings |
|---|---|---|---|---|---|
| Dependency | Advisory | Warning | High and critical | Baseline and remediate; do not block solely because they pre-date the baseline | Block high and critical |
| Code | Advisory | Warning | High | Baseline and remediate; do not block solely because they pre-date the baseline | Block high |
| IaC | Advisory | Warning | High and critical | Baseline and remediate; do not block solely because they pre-date the baseline | Block high and critical |
| Container | Advisory | Warning | High and critical | Baseline and remediate; do not block solely because they pre-date the baseline | Block high and critical |

Snyk Code has no distinct critical threshold in this integration, so `critical` maps to `high` for Code. All severities are published; the threshold controls policy, not evidence collection.

An existing finding is recorded in an approved baseline. A new finding is absent from that baseline, including a remediated finding that is reintroduced. Before Release 3, baselines must use stable finding identities and changes must be reviewable and auditable.

The current workflow can block on any finding at the threshold but cannot distinguish new from existing findings. Therefore `fail_build: true` must not be enabled organisation-wide until baseline/delta enforcement is implemented and tested. Earlier opt-in requires service-owner and pipeline-owner approval.

## Ignored findings

An ignore is valid only when its exception records the finding, justification, compensating controls, risk owner, approver, approval/expiry dates and remediation or risk-acceptance statement. Expired ignores become active automatically. Scanner-side ignores without matching approval evidence are unapproved.

## Failure behavior

| Condition | Observe | Warn | Enforce |
|---|---|---|---|
| Finding meets threshold | Publish; do not block | Prominent warning and review | Block if new or its exception expired |
| Existing baselined finding | Publish and track | Warn against remediation plan | Do not block solely because it is baselined |
| Scanner/network unavailable | Record `scan_operational_error`; never claim clean | Unsuccessful, initially non-required check | Fail required policy check |
| Timeout | Record operational failure | Unsuccessful, initially non-required check | Fail required policy check |
| Malformed/missing SARIF | Do not publish invalid evidence; record failure | Unsuccessful, initially non-required check | Fail required policy check |
| Requested SARIF upload fails | Retain evidence; record `result_publication_failed` | Unsuccessful, initially non-required check | Fail required policy check |

Operational failures may be overridden only through the emergency exception route and must never become clean results.

## Timeouts

| Job | Limit |
|---|---:|
| Validation | 10 minutes |
| Dependency | 15 minutes |
| Code | 20 minutes |
| IaC | 15 minutes |
| Container | 20 minutes |
| Publication | 10 minutes |
| Policy/summary | 5 minutes each |

A timed-out selected scan requires a successful retest.

## Exceptions

| Exception | Required approvers | Maximum duration |
|---|---|---:|
| Critical finding | Service owner and OCTO/AppSec or delegated security risk owner | 7 days |
| High finding | Service owner and OCTO/AppSec or delegated security risk owner | 30 days |
| Medium finding | Service owner and delegated security risk owner | 90 days |
| Scanner, timeout or publication failure | Service owner and central pipeline owner; notify OCTO/AppSec | 24 hours |
| Organisation-wide policy change | Central pipeline owner and OCTO/AppSec | Until next policy review |

The change author must not be the sole approver. Exceptions do not renew automatically. No exception may expose scanner credentials to PR-controlled code or weaken the protected-environment boundary.

Emergency approval must record the incident/change, repository and release, reason, approvers, compensating controls, expiry and retest. The previous known-good immutable workflow SHA remains the rollback target. Rollback must not restore repository-level scanner credentials.

## Retesting

A successful retest is required after remediation, exception expiry, scanner/network recovery, timeout, malformed/missing SARIF, publication failure, material scanned-content change, or workflow/policy change. Evidence must record the immutable workflow revision, run URL, conclusion and policy summary.

## Rollout

### Release 1 — Observe

- Run all selected supported scans and publish validated results.
- Use `fail_build: false` and do not block solely on findings.
- Classify operational/publication errors honestly.
- Establish baselines and measure compatibility, duration, failure rate and false positives.

Promotion requires completed Stage 9 production regression, an approved protected environment, representative consumers, documented baselines and no unresolved credential-boundary issue.

### Release 2 — Warn

- Report threshold violations prominently.
- Require review of high and critical findings.
- Exercise exception creation, expiry and retesting.
- Run the policy check on protected branches without making finding violations blocking.

Promotion requires baseline/delta capability, tested exception automation and rollback, agreed operational service levels, independent security approval and acceptable warning accuracy.

### Release 3 — Enforce

- Make `security / Policy result` the required branch-protection check.
- Block new findings at this policy's thresholds and expired/unapproved exceptions.
- Fail closed for operational errors, timeouts, malformed evidence and requested-publication failure.
- Maintain and periodically test emergency rollback.

## Initial decision and success gate

The first production release is **Release 1 — Observe**:

- report all severities;
- set `severity_threshold: high` as the future enforcement threshold;
- set `fail_build: false`;
- keep SARIF publication enabled;
- preserve fail-closed classification for validation, scanner and publication errors;
- do not activate organisation-wide finding enforcement until Release 3 gates pass.

Stage 10 is approved when service ownership and OCTO/AppSec approve this policy, exception ownership/evidence storage are named, rollout gates are accepted, workflow exit codes match the selected release, and branch protection matches the approved policy-check design.
