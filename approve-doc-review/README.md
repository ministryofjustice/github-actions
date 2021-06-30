# Approve documentation review

A Github Action to auto approve pull requests that contain only document
reviews.

The Cloud Platform team have a document reviewal process that requires
us to create a PR and have it reviewed by a colleague. This is cumbersome 
and unnecessary. An example of the PR in question looks like the following:

```
https://github.com/ministryofjustice/cloud-platform-user-guide/commit/de39ec0fd5f0bf97cc3e1054666e7709d56680e7
```

As you can see this is a fairly futile effort and time could be spent elsewhere.
This GitHub action checks to see if the only changes contain the words `last_review_on`
and whether the PR owner is from the Cloud Platform team. If the criteria is met the PR
is approved automatically.

To run this GitHub action you must create a GitHub personal access token and store it in a GitHub secret
in the repository in question. An example of the Action would look like:

```
name: Auto-approve a pull request

on:
  pull_request

env:
  PR_OWNER: ${{ github.event.pull_request.user.login }}
  GITHUB_OAUTH_TOKEN: ${{ secrets.DOCUMENT_REVIEW_GITHUB }}

jobs:
  check-diff:
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        os: [ubuntu-latest]

    steps:
      - name: Checkout PR code
        uses: actions/checkout@v2
      - run: |
          git fetch --no-tags --prune --depth=1 origin +refs/heads/*:refs/remotes/origin/*
      - name: owner of PR
        run: |
          echo "$PR_OWNER"
      - name: Run git diff against repository
        run: |
          git diff origin/main HEAD > changes
      - name: Auto-approval check
        id: approve_pr_check
        uses: ministryofjustice/github-actions/approve-doc-review@add-docs-check
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Approving PR
        uses: hmarr/auto-approve-action@v2

        if: steps.approve_pr_check.outputs.review_pr == 'true'
        with:
          github-token: "${{ secrets.GITHUB_TOKEN  }}"
```

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

NB: You do have to duplicate the conditional as shown. Although the
github documentation states that you can put the conditional at the
job level, that doesn't work, in this case.
