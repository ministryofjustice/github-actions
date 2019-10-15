# Branch deleter

A Github Action to delete a branch after its PR is merged.

People don't always remember to delete their branches on merge, so
this github action does it for you. The action only deletes a branch
when a PR is merged, not if the PR is closed without merging.

## USAGE

Create a file in your repo called `.github/workflows/delete-branch-after-merge.yml` with the
following contents:

```
on:
  pull_request:
    types: [closed]

jobs:
  delete-branch-after-merge:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
        if: github.event.pull_request.merged == true
      - uses: ministryofjustice/github-actions/branch-deleter@master
        if: github.event.pull_request.merged == true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

NB: You do have to duplicate the conditional as shown. Although the
github documentation states that you can put the conditional at the
job level, that doesn't work, in this case.
