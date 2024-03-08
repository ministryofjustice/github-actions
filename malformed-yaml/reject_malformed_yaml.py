
import yaml
from github_service import GitHubService as github_service


def get_changed_files_from_pr():
    return [
        "malformed-yaml/invalid.yaml",
    ]


def malformed_yaml(changed_files: list) -> bool:
    for file in changed_files:
        with open(file, "r") as f:
            try:
                yaml.safe_load(f)
            except yaml.YAMLError:
                return True
    return False


def get_github_token():
    token = "1234567890"
    return token


def does_pr_contain_malformed_yaml() -> bool:
    github = github_service(get_github_token())
    changed_files = github.get_changed_files_from_pr()
    if malformed_yaml(changed_files):
        return False
    return True


if __name__ == "__main__":
    does_pr_contain_malformed_yaml()
