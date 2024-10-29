# Check Version Pinning GitHub Action

This GitHub Action scans your workflow files to ensure all GitHub Actions are securely pinned to a SHA hash, rather than a version tag (`@v`). Using SHA pinning aligns with best practices to protect against unintended changes in third-party actions.

## Purpose

According to GitHub's security guidance, third-party actions should be pinned to a commit hash rather than a version tag for enhanced security. For instance, prefer this format:
```yaml
uses: oxsecurity/megalinter/flavors/python@32c1b3827a334c80026c654f31ee1b4801ad8798
```
over:
```yaml
uses: oxsecurity/megalinter/flavors/python@v1
```

This Action inspects workflows to detect and report any actions that are not SHA-pinned, helping to secure your CI/CD pipeline.

## Features

- Simple SHA Check: This Action scans workflows based on the string after the @ symbol to verify SHA pinning.

- Targeted Organisations: No organisations are treated as implicitly trusted, ensuring that all third-party actions must be SHA-pinned without exceptions.

- Customisable Scanning Modes: Run a full scan of your repository or focus on changes within a pull request.


## Inputs

`workflow_directory`

Specifies the directory to scan for workflow files. Defaults to .github/workflows if not set.

`scan_mode`

Defines the scope of the scan:

- full: Scans all workflows in the specified directory.
- pr_changes: Scans only changes within a pull request (PR).

## Outputs

Provides a list of any unpinned actions detected in the repository.

## Example Usage

Here's a typical workflow setup that uses this Action to enforce SHA pinning on actions:
```yaml
name: 🧪 Check Version Pinning

on:
  push:
    branches:
      - main
  pull_request:
  workflow_dispatch:

jobs:
  check-version-pinning:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
        with:
          fetch-depth: 0  # Disable shallow clones for a more comprehensive scan

      - name: Check for unpinned Actions
        uses: ministryofjustice/github-actions/check-version-pinning@ccf9e3a4a828df1ec741f6c8e6ed9d0acaef3490 # v18.5.0
        with:
          workflow_directory: ".github/workflows" # Or wherever your workflows are stored
          scan_mode: "full"  # or "pr_changes" for PR-specific scans
```

## Why This Action?

We initially considered using actionlint but found it too restrictive for our use case. This Action is lightweight and focuses solely on verifying SHA pinning for third-party actions, making it simpler and more tailored to specific security needs.

## Notes

This Action will:

- Flag any action with a version tag (e.g., @v1) rather than a SHA.

- Not detect cases where third-party actions do not use semantic versioning or the v prefix in version tags.

- Require all actions to be SHA-pinned, without any implicit trust for specific organisations like ministryofjustice or actions.

By adding this Action to your workflows, you can ensure a more secure CI/CD setup that prevents accidental usage of unpinned or untrusted actions.
