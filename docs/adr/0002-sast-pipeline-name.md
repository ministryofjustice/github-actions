# ADR 0002: Name the service SAST Pipeline

- Status: Accepted
- Date: 2026-08-30

## Decision

Use **SAST Pipeline** for the user-facing name and `sast.yml` for the reusable workflow.

## Rationale and consequences

The name describes the security capability rather than a vendor and covers the coordinated dependency, code, IaC and container checks. Documentation must define the included scan types because the term SAST is sometimes used only for source analysis.
