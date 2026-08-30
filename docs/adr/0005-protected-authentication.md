# ADR 0005: Use protected non-personal authentication

- Status: Accepted in principle; central provisioning pending
- Date: 2026-08-30

## Decision

Use a least-privilege non-personal OAuth client from an approved protected central environment. Do not store scanner credentials at consumer-repository level or expose them to PR-controlled code. A static non-personal token is an approved fallback only.

## Rationale and consequences

Controlled testing proved repository secrets were available to same-repository PR workflow code. Protected execution prevents that boundary failure and supports central rotation/revocation. Central environment ownership and cross-repository behavior must be proven before release; missing provisioning must not be bypassed with `secrets: inherit`.
