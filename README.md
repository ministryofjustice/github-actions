# 💣 GitHub Actions

[![Ministry of Justice Repository Compliance Badge](https://github-community.service.justice.gov.uk/repository-standards/api/github-actions/badge)](https://github-community.service.justice.gov.uk/repository-standards/github-actions)

> A single pane of glass for discovering **reusable GitHub Actions and Workflows** maintained across the Ministry of Justice.

This repository is primarily a **discovery hub** for teams across MoJ to advertise their reusable GitHub automation tools. It also hosts centrally owned workflows approved for organisation-wide use, beginning with the [SAST Pipeline](docs/sast-pipeline.md).

## 🧭 Purpose of This Repository

This repo acts as a **directory** of reusable GitHub Actions and Workflows maintained by MoJ teams. It is designed to:

- Encourage **reuse** of existing automation
- Avoid duplication of effort
- Promote **visibility** of automation best practices
- Help teams **contribute and maintain** shared solutions

## 📚 Repositories Hosting Reusable Actions & Workflows

Below is a growing list of repositories that contain reusable Actions or Workflows:

| Repository                                                                                                          | Maintainer             | Description                                                                                                                                         |
| ------------------------------------------------------------------------------------------------------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| [analytical-platform-github-actions](https://github.com/ministryofjustice/analytical-platform-github-actions)       | Analytical Platform    | Analytical Platform GitHub Actions                                                                                                                  |
| [hmpps-github-actions](https://github.com/ministryofjustice/hmpps-github-actions)                                   | HMPPS                  | Github actions for HMPPS projects                                                                                                                   |
| [laa-reusable-github-actions](https://github.com/ministryofjustice/laa-reusable-github-actions)                     | LAA                    | A collection of re-useable GitHub actions                                                                                                           |
| [modernisation-platform-github-actions](https://github.com/ministryofjustice/modernisation-platform-github-actions) | Modernisation Platform | A collection of reusable GitHub Actions for the Modernisation Platform, designed to streamline and enhance workflows across our projects            |
| [opg-github-actions](https://github.com/ministryofjustice/opg-github-actions)                                       | OPG                    | OPG shared GitHub composite actions for workflows                                                                                                   |
| [cloud-platform-github-workflows](https://github.com/ministryofjustice/cloud-platform-github-workflows)             | Cloud Platform         | A central place to store Cloud Platforms reusable workflows.                                                                                        |
| [data-engineering-github-workflows](https://github.com/ministryofjustice/data-engineering-github-actions)           | Data Engineering       | A central place to store Data Engineering reusable workflows.                                                                                       |
| [devsecops-actions](https://github.com/ministryofjustice/devsecops-actions)                                         | PandA                  | ⚡️A collection of reusable GitHub Actions that standardise DevSecOps security scanning i.e. SCA, SAST, DAST, secrets, IaC, and container security.  |
| _Add yours here_                                                                                                    | You?                   | Open a PR to add your repository and reusable components                                                                                            |

> ✨ Want to list your repository? [See how to get involved](#-how-to-contribute)

## 🔍 What Are Reusable GitHub Actions?

GitHub Actions can automate workflows for CI/CD, security checks, infrastructure provisioning, and more.

There are two main types of reusable automation:

- **Reusable Actions**: Individual building blocks that perform a specific task (e.g. `setup-terraform`, `slack-notify`)
- **Reusable Workflows**: Complete pipelines composed of multiple steps, which can be invoked using `workflow_call`

**Learn more:**

- [GitHub Docs: Reusing workflows](https://docs.github.com/en/actions/how-tos/sharing-automations/reuse-workflows)
- [GitHub Docs: Creating actions](https://docs.github.com/en/actions/how-tos/sharing-automations/creating-actions)

## 🧪 How to Use a Reusable Workflow

### Analytical Platform

You can call reusable workflows from by using the `uses` attribute. Example:

```yaml
jobs:
  zizmor:
    name: Zizmor
    permissions:
      actions: read
      contents: read
      security-events: write
    uses: ministryofjustice/analytical-platform-github-actions/.github/workflows/reusable-zizmor.yml@<commit SHA> # <version>
```

### SCA

To use the MoJ's central SCA (Software Composition Analysis) composite action in your repository, use the basic snippet:

```yaml
name: SCA
run-name: SCA ⚡️

on:
  schedule:
    - cron: "0 0 * * *" # Run daily at midnight UTC

  pull_request:
    branches: ["main"]

permissions: {}

jobs:
  sca:
    name: Software Composition Analysis
    runs-on: ubuntu-latest

    permissions:
      contents: write
      pull-requests: write
      issues: write
      security-events: write

    steps:
      - name: Run SCA
        uses: ministryofjustice/devsecops-actions/sca@v1.3.0
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

### SAST Pipeline

The centrally owned SAST Pipeline provides dependency, code, infrastructure-as-code and existing-container scanning with consistent SARIF publication and failure behaviour.

See the [consumer contract and pinned example](docs/sast-pipeline.md). Consumers should reference a reviewed release tag or complete commit SHA, never a mutable branch.

Additional documentation:

- [Maintainer guide](docs/sast-maintainers.md)
- [Migration from the LAA workflow](docs/sast-migration.md)
- [Enforcement and exception policy](docs/sast-enforcement-policy.md)
- [Architecture decisions](docs/adr/README.md)

> Make sure to check each repository’s README for usage instructions and required inputs/secrets.

## 🪝 Git hook

Repository uses [MoJ DevSecOps hooks](https://github.com/ministryofjustice/devsecops-hooks) to ensure `pre-commit` git hook is evaluated for series of checks before pushing the changes from staging area. Engineers should ensure `pre-commit` hook is configured and activated.

1. **Installation**:

   Ensure [prek](https://github.com/j178/prek?tab=readme-ov-file#installation) is installed globally

   Linux / MacOS

   ```bash
   curl --proto '=https' --tlsv1.2 -LsSf https://raw.githubusercontent.com/ministryofjustice/devsecops-hooks/e85ca6127808ef407bc1e8ff21efed0bbd32bb1a/prek/prek-installer.sh | sh
   ```

   or

   ```bash
   brew install prek
   ```

   Windows

   ```bash
   powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/ministryofjustice/devsecops-hooks/e85ca6127808ef407bc1e8ff21efed0bbd32bb1a/prek/prek-installer.ps1 | iex"
   ```

2. **Activation**

   Execute the following command in the repository directory

   ```bash
   prek install
   ```

3. **Test**

   To dry-run the hook

   ```bash
   prek run
   ```

## 🔧 Configuration

### Exclusion list

One can exclude files and directories by adding them to `exclude` property. Exclude property accepts [regular expression](https://pre-commit.com/#regular-expressions).

Ignore everything under `reports` and `docs` directories for `baseline` hook as an example.

```yaml
   repos:
     - repo: https://github.com/ministryofjustice/devsecops-hooks
       rev: v1.0.0
       hooks:
         - id: baseline
            exclude: |
            ^reports/|
            ^docs/
```

Or one can also create a file with list of exclusions.

```yaml
repos:
  - repo: https://github.com/ministryofjustice/devsecops-hooks
    rev: v1.0.0
    hooks:
      - id: baseline
        exclude: .pre-commit-ignore
```

## 🤝 How to Contribute

We welcome contributions from all MoJ teams! Here’s how to get involved:

1. **Host your own GitHub Actions or Workflows** in your team’s repo.
2. **Document** their purpose, usage, and versioning clearly in the `README.md`.
3. **Open a pull request** to this repository adding your repo to the directory table above.
4. Optionally, provide example usages to help others adopt them faster.

> 📌 For new teams: We recommend versioning your actions using tags (e.g. `v1`, `v2.1.0`) to ensure stability.
