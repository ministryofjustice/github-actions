# ADR 0006: Use immutable version references

- Status: Accepted
- Date: 2026-08-30

## Decision

Use a full commit SHA for pre-release and high-assurance callers. Publish reviewed immutable semantic-version tags for supported releases. Never recommend mutable branches in production.

## Rationale and consequences

SHAs provide exact provenance; tags provide managed discoverability. Published tags must not move. Breaking contract changes require a major version, and rollback uses the previous known-good SHA. Tag and rollback tests are release gates.
