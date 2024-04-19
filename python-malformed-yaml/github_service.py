from github import Github


class GitHubService():
    """A class to interact with GitHub API."""

    def __init__(self, github_token, repository_name: str, pr_number: int):
        self.github_token = github_token
        self.pr_number = pr_number
        self.repository_name = repository_name
        self.client = Github(self.github_token)

    def get_changed_files_from_pr(self) -> list:         
        return [file.filename for file in self.client.get_repo(
            self.repository_name).get_pull(self.pr_number).get_files()]

    def fail_pr(self, message: str):
        self.client.get_repo(
            self.repository_name).get_pull(self.pr_number).create_issue_comment(
                message)
        return 1