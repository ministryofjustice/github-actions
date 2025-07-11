# 💣 GitHub Actions

[![Ministry of Justice Repository Compliance Badge](https://github-community.service.justice.gov.uk/repository-standards/api/github-actions/badge)](https://github-community.service.justice.gov.uk/repository-standards/github-actions)

> A single pane of glass for discovering **reusable GitHub Actions and Workflows** maintained across the Ministry of Justice.

This repository does **not** host any GitHub Actions or Workflows itself. Instead, it serves as a **discovery hub** for teams across MoJ to advertise their reusable GitHub automation tools, enabling greater collaboration, consistency, and efficiency across projects.

---

## 🧭 Purpose of This Repository

This repo acts as a **directory** of reusable GitHub Actions and Workflows maintained by MoJ teams. It is designed to:

- Encourage **reuse** of existing automation
- Avoid duplication of effort
- Promote **visibility** of automation best practices
- Help teams **contribute and maintain** shared solutions

---

## 📚 Repositories Hosting Reusable Actions & Workflows

Below is a growing list of repositories that contain reusable Actions or Workflows:

| Repository                                                                                                          | Maintainer             | Description                                                                                                                              |
| ------------------------------------------------------------------------------------------------------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| [analytical-platform-github-actions](https://github.com/ministryofjustice/analytical-platform-github-actions)       | Analytical Platform    | Analytical Platform GitHub Actions                                                                                                       |
| [hmpps-github-actions](https://github.com/ministryofjustice/hmpps-github-actions)                                   | HMPPS                  | Github actions for HMPPS projects                                                                                                        |
| [laa-reusable-github-actions](https://github.com/ministryofjustice/laa-reusable-github-actions)                     | LAA                    | A collection of re-useable GitHub actions                                                                                                |
| [modernisation-platform-github-actions](https://github.com/ministryofjustice/modernisation-platform-github-actions) | Modernisation Platform | A collection of reusable GitHub Actions for the Modernisation Platform, designed to streamline and enhance workflows across our projects |
| [opg-github-actions](https://github.com/ministryofjustice/opg-github-actions)                                       | OPG                    | OPG shared GitHub composite actions for workflows                                                                                        |
| _Add yours here_                                                                                                    | You?                   | Open a PR to add your repository and reusable components                                                                                 |

> ✨ Want to list your repository? [See how to get involved](#-how-to-contribute)

---

## 🔍 What Are Reusable GitHub Actions?

GitHub Actions can automate workflows for CI/CD, security checks, infrastructure provisioning, and more.

There are two main types of reusable automation:

- **Reusable Actions**: Individual building blocks that perform a specific task (e.g. `setup-terraform`, `slack-notify`)
- **Reusable Workflows**: Complete pipelines composed of multiple steps, which can be invoked using `workflow_call`

**Learn more:**

- [GitHub Docs: Reusing workflows](https://docs.github.com/en/actions/how-tos/sharing-automations/reuse-workflows)
- [GitHub Docs: Creating actions](https://docs.github.com/en/actions/how-tos/sharing-automations/creating-actions)

## 🧪 How to Use a Reusable Workflow

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

> Make sure to check each repository’s README for usage instructions and required inputs/secrets.

---

## 🤝 How to Contribute

We welcome contributions from all MoJ teams! Here’s how to get involved:

1. **Host your own GitHub Actions or Workflows** in your team’s repo.
2. **Document** their purpose, usage, and versioning clearly in the `README.md`.
3. **Open a pull request** to this repository adding your repo to the directory table above.
4. Optionally, provide example usages to help others adopt them faster.

> 📌 For new teams: We recommend versioning your actions using tags (e.g. `v1`, `v2.1.0`) to ensure stability.
