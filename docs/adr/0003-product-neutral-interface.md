# ADR 0003: Keep the interface product-neutral

- Status: Accepted
- Date: 2026-08-30

## Decision

Expose scan selection, paths, severity, publication and enforcement concepts without Snyk-specific input names.

## Rationale and consequences

Consumers depend on the capability rather than a vendor contract, reducing migration cost if implementation changes. Some implementation limitations, such as Code severity mapping, still require documentation. Vendor-specific administration remains internal to maintainers.
