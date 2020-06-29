# Sync gh-pages branch

GitHub only allows you to publish your "master" or "gh-pages" branch as a GitHub Pages website.

This is a problem if you want to change your repo to use "main" instead of "master" as the default branch.

This action automatically keeps the "gh-pages" branch up to date with respecto the "main" branch, so you can publish the "gh-pages" branch via GitHub Pages, but still work on "main" as your default branch.

## USAGE

Create a file in your repo called `.github/workflows/gh-pages-sync.yml` with the following contents:

```
name: Keep gh-pages in sync wrt. main

on:
  push:
    branches: [main]

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - run: git checkout -b gh-pages
      - run: git push origin gh-pages --force
```

## Limitations

* Any changes you make directly to the "gh-pages" branch will be lost, the next time someone updates "main"
* You cannot use this if you are publishing the "docs" folder of your master branch. It is currently a limitation of GitHub Pages that you can only publish the "docs" folder of the "master" branch - not any other branch.
