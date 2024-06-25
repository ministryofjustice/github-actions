import os
from pathlib import Path
import logging
import sys
import re
import yaml
from github_pull_request import GitHubPullRequest as github_pull_request


logging.basicConfig(
    format="%(asctime)s %(levelname)-8s %(message)s",
    level="INFO",
    datefmt="%Y-%m-%d %H:%M:%S",
)

logger = logging.getLogger(__name__)


def get_github_env() -> tuple[str, str, str]:
    """
    Function to collect the three required GitHub
    environment variables
    """
    token = os.getenv("GITHUB_TOKEN")
    pr_number = os.getenv("PR_NUMBER")
    repo = os.getenv("GITHUB_REPOSITORY")
    if not token:
        raise ValueError("No GITHUB_TOKEN.")
    if not pr_number:
        raise ValueError("No PR_NUMBER.")
    return token, repo, pr_number

def get_extant_files(files: list[str]) -> list[str]:
    """
    Check that a list of files exists and prune any that do not.
    Return extant files only.
    """
    return [file for file in files if Path(file).exists()]

def get_changed_yaml_files_from_pr() -> list[str]:
    """
    Collect a list of all the new or modified YAML files (with path)
    in a PR, excluding files in a 'secret/' directory.
    """
    token, repository_name, pr = get_github_env()
    github_pr = github_pull_request(token, repository_name, int(pr))
    # We assume there must always be some changed or new files in a PR
    changed_files = github_pr.get_changed_files_from_pr()
    yml_pattern = re.compile("\\.yml$|\\.yaml$")
    skip_pattern = re.compile("secret/")

    return [
        file for file in changed_files if yml_pattern.search(file) and not skip_pattern.search(file)
    ]

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
        with open(y, encoding="utf-8") as stream:
            try:
                for _ in yaml.safe_load_all(stream):
                    pass
            except yaml.YAMLError as exc:
                malformed_yaml_files_and_errors.append(f"\n{str(y)}:\n{str(exc)}")
    return malformed_yaml_files_and_errors

def malformed_yaml_files_message(files_and_errors: list):
    """
    Compose message to display in the PR.
    """
    msg = "😱 The following malformed YAML files and related errors were found:\n"
    msg += "\n".join(files_and_errors) + "\n\n🥺 Please correct them and resubmit this PR."
    return msg

def main():
    """
    Function to collect the new or modified YAML files from the PR that
    are malformed, report these to the user, and request changes.
    """
    token, repository_name, pr = get_github_env()
    github_pr = github_pull_request(token, repository_name, int(pr))


    changed_yaml_files = get_changed_yaml_files_from_pr()
    if not changed_yaml_files:
        logger.info("🫧 No new or modified YAML files.")
        return False


    malformed_yaml_files_and_errors = get_malformed_yaml_files_and_errors(
        get_extant_files(changed_yaml_files)
    )
    if malformed_yaml_files_and_errors:
        msg = malformed_yaml_files_message(malformed_yaml_files_and_errors)
        github_pr.fail_pr(message=msg)
        logger.error(msg)
        return True

    logger.info("🤩 PR YAML files all OK!")
    return False


if __name__ == "__main__":
    if main():
        sys.exit(1)
