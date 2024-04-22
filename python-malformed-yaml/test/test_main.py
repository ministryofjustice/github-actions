import os
from unittest.mock import patch
import unittest
from main import get_changed_yaml_files_from_pr, get_malformed_yaml_files_and_errors, main
from github_service import GitHubService as github_service


class TestMain(unittest.TestCase):


    @patch.object(github_service, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})
    def test_get_changed_yaml_files_from_pr_all_combos(self, mock_github_service):
        mock_github_service.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "b.yaml", "some/other/c.yml", "secret/d.yml", "secret/directory/e.yaml", "really/secret/stuff/f.txt"
        ]
        result = get_changed_yaml_files_from_pr()
        self.assertEqual(result, ["b.yaml", "some/other/c.yml"])

    @patch.object(github_service, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})
    def test_get_changed_yaml_files_from_pr_no_yaml(self, mock_github_service):
        mock_github_service.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "secret/non/yaml.csv", "some.yml/but/not/a/yml.txt"
        ]
        result = get_changed_yaml_files_from_pr()
        self.assertEqual(result, [])