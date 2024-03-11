import os

import yaml
from github_service import GitHubService as github_service


def malformed_yaml(changed_files: list) -> set:
    files = set()
    for file in changed_files:
        with open(file, "r") as f:
            try:
                yaml.safe_load(f)
            except yaml.YAMLError:
                files.add(file)
    return files


def get_github_env() -> tuple:
    pr_number = os.getenv("PR_NUMBER")
    token = os.getenv("GITHUB_TOKEN")
    repo = os.getenv("REPOSITORY_NAME")
    if not token or not pr_number or not repo:
        raise ValueError("No GITHUB_TOKEN or PR_NUMBER env var found. Please make this available via the github actions workflow\nhttps://help.github.com/en/articles/virtual-environments-for-github-actions#github_token-secret.")
    return token, repo, pr_number


def message(files: set):
    msg = "😱 The following files contain malformed YAML:\n -"
    msg += "\n -".join(files)
    return msg


def does_pr_contain_malformed_yaml() -> bool:
    token, repository_name, pr = get_github_env()
    github = github_service(token, repository_name, pr)

    changed_files = github.get_changed_files_from_pr()
    malformed_files = malformed_yaml(changed_files)
    if malformed_files:
        github.fail_pr(message(malformed_files))
        print(f"PR contains malformed YAML: {malformed_files}")
        return True
    print("PR does not contain malformed YAML")
    return False


if __name__ == "__main__":
    if does_pr_contain_malformed_yaml():
        # Lets exit with a non-zero status code to indicate failure to GitHub Actions
        os._exit(1)
