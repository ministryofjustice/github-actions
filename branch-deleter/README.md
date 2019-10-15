# Branch deleter

A Github Action to delete a branch after its PR is merged.

People don't always remember to delete their branches on merge, so
this github action does it for you. The action only deletes a branch
when a PR is merged, not if the PR is closed without merging.

NB: This function is now available as an option on Github.
Look for the "Merge button" section in [settings](https://github.com/ministryofjustice/github-actions/settings)
and turn on "Automatically delete head branches." This code is left
here for reference only.

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

`GITHUB_TOKEN` is provided automatically by github actions. You do
not need to do anything extra to make it available. Just use the
content above, exactly as shown.

NB: You do have to duplicate the conditional as shown. Although the
github documentation states that you can put the conditional at the
job level, that doesn't work, in this case.
