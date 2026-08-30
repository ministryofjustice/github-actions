# ADR 0007: Introduce enforcement in phases

- Status: Accepted for Release 1
- Date: 2026-08-30

## Decision

Start with Observe (`severity_threshold: high`, `fail_build: false`), progress through Warn, and enforce only after baseline/delta, exception and rollback controls are tested. The future stable required check is `security / Policy result`.

## Rationale and consequences

Immediate organisation-wide blocking would make existing findings and compatibility problems delivery failures before accuracy is known. Operational and requested-publication errors still fail honestly. Release 3 will block new high/critical findings according to the approved enforcement policy.
