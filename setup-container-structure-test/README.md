# Set Up Container Structure Test

This action installs Google's [Container Structure Test](https://github.com/GoogleContainerTools/container-structure-test) tool.

> [!warning]
> Container Structure Test is not an officially supported Google project, and is currently in maintainence mode.
> However, releases are still being created.

## Usage

> [!warning]
> This action only works with versions of Container Structure Test that are released to GitHub.
> v1.17.0 onwards.

```yaml
jobs:
  container-structure-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        id: checkout
        uses: actions/checkout@v4

      - name: Set Up Container Structure Test
        id: setup_container_structure_test
        uses: ministryofjustice/github-actions/setup-container-structure-test@main

      - name: Run Container Structure Test
        id: run_container_structure_test
        run: |
          container-structure-test ...
```

Specifying a version

```yaml
      - name: Set Up Container Structure Test
        id: setup_container_structure_test
        uses: ministryofjustice/github-actions/setup-container-structure-test@main
        with:
          version: v1.17.0
```