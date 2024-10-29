# Check Version Pinning GitHub Action

This Action scans your workflow files for untrusted GitHub Actions that are pinned to a version (`@v`) rather than a SHA hash.

## Inputs

### `workflow_directory`
The directory to scan for workflow files. Default is `.github/workflows`.

### `scan_mode`
The type of scan you wish to undertake:
- full = the whole repository.
- pr_changes = only changes in a pr.

## Outputs

### `found_unpinned_actions`
A boolean indicating if any unpinned actions were found.

## Example usage
```yaml
jobs:
  check-version-pinning:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check for unpinned Actions
        uses: ministryofjustice/check-version-pinning-action@v1
        with:
          workflow_directory: ".github/workflows"
          scan_mode: "pr_changes"  # Use "full" for a full repo scan
```
