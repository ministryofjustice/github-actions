

class GitHubService():
    """A class to interact with GitHub API."""

    def __init__(self, github_token):
        self.github_token = github_token

    def get_changed_files_from_pr(self):
        return [
            "malformed-yaml/invalid.yaml",
        ]
