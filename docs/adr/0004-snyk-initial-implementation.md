# ADR 0004: Use Snyk as the initial implementation

- Status: Accepted for the initial release
- Date: 2026-08-30

## Decision

Use a pinned Snyk CLI for dependency, code, IaC and container scanning while retaining a product-neutral caller contract.

## Rationale and consequences

Snyk already supports the source LAA workflow, the required scan types and SARIF publication. The choice requires organisation entitlement, non-personal credentials, CLI supply-chain review and explicit handling of Snyk Code's severity model. Replacing Snyk must preserve or version the public contract.
