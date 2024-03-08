from github import Github


class GitHubService():
    """A class to interact with GitHub API."""

    def __init__(self, github_token, pr_number: int):
        self.github_token = github_token
        self.pr_number = pr_number
        self.client = Github(self.github_token)

    def get_changed_files_from_pr(self):
        self.client.
        return [
            "malformed.yaml",
            "formed.yaml",
        ]

    def fail_pr(self, message: str):
        return 1
