# Clean Actions Runner

> [!CAUTION]
> The steps performed in this action are destructive, please review before using!

Standard Github-hosted runner are only guaranteed 14GB of SSD storage ([source](https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners#standard-github-hosted-runners-for-public-repositories)), this action removes bundled software as discussed in [this](https://github.com/actions/runner-images/issues/2840) GitHub issue.

This action is useful when working with large container images.

## Usage

To run the cleanup operation, you will need to explicitly set `confirm: true`, for example:

```yaml
- name: Clean Actions Runner
  id: clean_actions_runner
  uses: ministryofjustice/github-actions/clean-actions-runner@main
  with:
    confirm: true
```

To retain a specific piece of software, set its input to `false`, for example:

```yaml
- name: Clean Actions Runner
  id: clean_actions_runner
  uses: ministryofjustice/github-actions/clean-actions-runner@main
  with:
    confirm: true
    remove_opt_hostedtoolcache: false
```

Using the default options should reclaim about 29GB.

## Inputs

| Input | Default | Required | Description |
|:---|:---:|:---:|:---|
| `confirm` | `false` | `true`   | Confirm that you want to remove the software |
| `remove_opt_google` | `true` | `true` | Remove /opt/google (347MB) |
| `remove_opt_hostedtoolcache` | `true` | `true` | Remove /opt/hostedtoolcache (8.5GB) |
| `remove_opt_microsoft` | `true` | `true` | Remove /opt/microsoft (743MB) |
| `remove_opt_pipx` | `true` | `true` | Remove /opt/pipx (437MB) |
| `remove_usr_lib_firefox` | `true` | `true` | Remove /usr/lib/firefox (257MB) |
| `remove_usr_lib_google_cloud_sdk` | `true` | `true` | Remove /usr/lib/google-cloud-sdk (916MB) |
| `remove_usr_lib_heroku` | `true` | `true` | Remove /usr/lib/heroku (280MB) |
| `remove_usr_lib_llvm_13` | `true` | `true` | Remove /usr/lib/llvm-13 (448MB) |
| `remove_usr_lib_llvm_14` | `true` | `true` | Remove /usr/lib/llvm-14 (486MB) |
| `remove_usr_lib_llvm_15` | `true` | `true` | Remove /usr/lib/llvm-15 (514MB) |
| `remove_usr_lib_mono` | `true` | `true` | Remove /usr/lib/mono (423MB) |
| `remove_usr_local_julia` | `true` | `true` | Remove /usr/local/julia* (856MB) |
| `remove_usr_local_lib_android` | `true` | `true` | Remove /usr/local/lib/android (7.6GB) |
| `remove_usr_local_lib_node_modules` | `true` | `true` | Remove /usr/local/lib/node_modules (1.1GB) |
| `remove_usr_local_share_chromium` | `true` | `true` | Remove /usr/local/share/chromium (542MB) |
| `remove_usr_local_share_powershell` | `true` | `true` | Remove /usr/local/share/powershell (1.2GB) |
| `remove_usr_share_dotnet` | `true` | `true` | Remove /usr/share/dotnet (1.6GB) |
| `remove_usr_share_miniconda` | `true` | `true` | Remove /usr/share/miniconda (658MB) |
| `remove_usr_share_swift` | `true` | `true` | Remove /usr/share/swift (2.6GB) |
