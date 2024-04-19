import os
from pathlib import Path
import re
import yaml
from github_service import GitHubService as github_service


def get_github_env() -> tuple:
    token = os.getenv("GITHUB_TOKEN")
    pr_number = os.getenv("PR_NUMBER")
    repo = os.getenv("REPOSITORY_NAME")
    if not token:
        raise ValueError("No GITHUB_TOKEN.")
    if not pr_number:
        raise ValueError("No PR_NUMBER.")
    if not repo:
        raise ValueError("No REPOSITORY_NAME.")
    return token, repo, pr_number


def get_changed_yaml_files_from_pr() -> list:
    token, repository_name, pr = get_github_env()
    github = github_service(token, repository_name, int(pr))
    changed_files = github.get_changed_files_from_pr()
    pattern = re.compile("\\.yml$|\\.yaml$")
    changed_yaml_files = [file for file in changed_files if pattern.search(file)]
    return changed_yaml_files

def get_malformed_yaml_files(yaml_files: list) -> list:
    print(f"WD: {os.getcwd()}")
    # p = Path.cwd()
    # os.chdir(p.parent)
    # print(f"NEW WD: {os.getcwd()}")

    malformed_yaml_files = []
    for y in yaml_files:
        with open(y) as stream:
            try:
                yaml.safe_load(stream)
            except yaml.YAMLError as exc:
                malformed_yaml_files.append(f"\n{str(y)}:\n{str(exc)}")   
    return malformed_yaml_files     

def main():
    changed_yaml_files = get_changed_yaml_files_from_pr()
    print(changed_yaml_files)
    malformed_yaml_files = get_malformed_yaml_files(changed_yaml_files)
    print(malformed_yaml_files)

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
