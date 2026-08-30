# ADR 0001: Host SAST Pipeline in `github-actions`

- Status: Accepted
- Date: 2026-08-30

## Decision

Host the organisation-facing reusable workflow in `ministryofjustice/github-actions`.

## Rationale and consequences

The repository is the MoJ discovery and governance location for shared GitHub automation and has central CODEOWNERS. This avoids making an LAA-owned repository the long-term organisation contract. Central ownership increases the review and release burden and requires coordination with the former LAA maintainer.
