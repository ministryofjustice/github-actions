# Rbac permissions check

A GitHub Action to check if a user is allowed to modify a [cloud-platform-environments](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/getting-started/env-create.html#namespace-yaml-files) namespace.

When a user creates a pull request on the [cloud-platform-environments](https://github.com/ministryofjustice/cloud-platform-environments/) repository this check is triggered confirming the users access. If the user is in the relevant github team specified in the namespaces [rbac file](https://github.com/ministryofjustice/cloud-platform-environments/blob/main/namespaces/live-1.cloud-platform.service.justice.gov.uk/abundant-namespace-dev/01-rbac.yaml) the check will pass, However if the user isn't a member of the team the Action will fail and a message will be posted into the comments of the PR.

This ensures only team members can modify, create and delete namespaces they have ownership of.
