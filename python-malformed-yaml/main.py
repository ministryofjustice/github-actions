import os
import logging
import re
import yaml
from github_service import GitHubService as github_service


logging.basicConfig(
    format="%(asctime)s %(levelname)-8s %(message)s",
    level="INFO",
    datefmt="%Y-%m-%d %H:%M:%S",
)

logger = logging.getLogger(__name__)


def get_github_env() -> tuple[str, str, str]:
    """
    Function to collect the three required GitHub
    environmnet variables
    """
    token = os.getenv("GITHUB_TOKEN")
    pr_number = os.getenv("PR_NUMBER")
    repo = os.getenv("GITHUB_REPOSITORY")
    if not token:
        raise ValueError("No GITHUB_TOKEN.")
    if not pr_number:
        raise ValueError("No PR_NUMBER.")
    return token, repo, pr_number


def get_changed_yaml_files_from_pr() -> list[str]:
    """
    Collect a list of all the new or modified YAML files
    in a PR, except those in a 'secret/' directory.
    """
    token, repository_name, pr = get_github_env()
    github = github_service(token, repository_name, int(pr))
    changed_files = github.get_changed_files_from_pr()
    pattern = re.compile("\\.yml$|\\.yaml$")
    skip_pattern = re.compile("secret/")
    changed_yaml_files = [
        file for file in changed_files if pattern.search(file) and not skip_pattern.search(file)
    ]
    return changed_yaml_files

def get_malformed_yaml_files_and_errors(yaml_files: list[str]) -> list[str]:
    """
    Input:
        yaml_files: List of YAML files to be tested for correct format
    Output:
        malformed_yaml_files: List of those YAML files that are malformed and
        their error messages.
    """
    malformed_yaml_files_and_errors= []
    for y in yaml_files:
        with open(y) as stream:
            try:
                yaml.safe_load(stream)
            except yaml.YAMLError as exc:
                malformed_yaml_files_and_errors.append(f"\n{str(y)}:\n{str(exc)}")
    return malformed_yaml_files_and_errors

def malformed_yaml_files_message(files_and_errors: list):
    """
    Compose message to display in the PR.
    """
    msg = "😱 The following malformed YAML files and related errors were found:\n"
    msg += "\n".join(files_and_errors)
    return msg

def main():
    """
    Function to collect the new or modified YAML files from the PR that
    are malformed, report these to the user, and request changes.
    """
    token, repository_name, pr = get_github_env()
    github = github_service(token, repository_name, int(pr))

    changed_yaml_files = get_changed_yaml_files_from_pr()
    malformed_yaml_files = get_malformed_yaml_files_and_errors(changed_yaml_files)
    if malformed_yaml_files:
        msg = malformed_yaml_files_message(malformed_yaml_files)
        github.fail_pr(message=msg)
        logger.error(msg)
    else:
        logger.info("PR YAML files all OK!")


if __name__ == "__main__":
    main()
