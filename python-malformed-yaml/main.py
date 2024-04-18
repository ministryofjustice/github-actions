import os
import yaml
import pathlib
from github_service import GitHubService as github_service


def get_github_env() -> tuple:
    pr_number = os.getenv("PR_NUMBER")
    token = os.getenv("GITHUB_TOKEN")
    repo = os.getenv("REPOSITORY_NAME")
    if not token or not pr_number or not repo:
        raise ValueError("No GITHUB_TOKEN or PR_NUMBER env var found. Please make this available via the github actions workflow\nhttps://help.github.com/en/articles/virtual-environments-for-github-actions#github_token-secret.")
    return token, repo, pr_number


def get_changed_yaml_files_from_pr():
    token, repository_name, pr = get_github_env()
    github = github_service(token, repository_name, pr)
    changed_files = github.get_changed_files_from_pr()
    print(changed_files)


def main():
    print("!!!OUTPUT!!!")
    get_changed_yaml_files_from_pr()

# def main():

#     yml_files = [p for p in pathlib.Path(".").rglob(
#         '*') if p.suffix in [".yml", ".yaml"]]
#     yml_files = [y for y in yml_files if "secret/" not in str(y)]

#     malformed_yaml = []
#     for y in yml_files:
#         with open(y) as stream:
#             try:
#                 yaml.safe_load(stream)
#             except yaml.YAMLError as exc:
#                 malformed_yaml.append(f"\n{str(y)}:\n{str(exc)}")

#     if malformed_yaml != []:
#         error_message = (
#             "Malformed YAML detected:\n" +
#             "\n".join(malformed_yaml) +
#             ("\n Please correct and resubmit this PR.")
#         )
#         raise Exception(error_message)
#     else:
#         print("All YAML files OK!")


if __name__ == "__main__":
    main()
