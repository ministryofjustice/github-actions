import os
import logging
from unittest.mock import patch, MagicMock
import unittest
from main import get_changed_yaml_files_from_pr, get_malformed_yaml_files_and_errors, main
from github_pull_request import GitHubPullRequest as github_pull_request


class TestMain(unittest.TestCase):


    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})
    def test_get_changed_yaml_files_from_pr_ignores_secret_directory(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "secret/d.yml", "secret/directory/e.yaml", "really/secret/stuff/f.yaml"
        ]
        result = get_changed_yaml_files_from_pr()
        self.assertEqual(result, [])

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})
    def test_get_changed_yaml_files_from_pr_ignores_non_yamls(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "b.rmd", "some/other/yml.txt", "not/a/yaml/file.py"
        ]
        result = get_changed_yaml_files_from_pr()
        self.assertEqual(result, [])

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})
    def test_get_changed_yaml_files_from_pr_returns_only_yaml(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "b.yaml", "some/other/c.yml", "another/yaml.yaml", "secret/file.yaml",
        ]
        result = get_changed_yaml_files_from_pr()
        self.assertEqual(result, ["b.yaml", "some/other/c.yml", "another/yaml.yaml"])

    def test_get_malformed_yaml_files_and_errors(self):
        result = get_malformed_yaml_files_and_errors(
            yaml_files=[
                "test/test_yaml_files/bad.yaml",
                "test/test_yaml_files/bad.yml",
                "test/test_yaml_files/good.yaml",
            ]
        )
        expected = [
            '\ntest/test_yaml_files/bad.yaml:\nwhile scanning ' +
            'a quoted scalar\n  in "test/test_yaml_files/' +
            'bad.yaml", line 2, column 14\nfound unexpected end of stream\n  in ' +
            '"test/test_yaml_files/bad.yaml", line 3, column 1',
            '\ntest/test_yaml_files/bad.yml:\nwhile scanning a' +
            ' quoted scalar\n  in "test/test_yaml_files/' +
            'bad.yml", line 2, column 6\nfound unexpected end of stream\n  in ' +
            '"test/test_yaml_files/bad.yml", line 3, column 1'
        ]
        self.assertEqual(expected, result)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_true(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
                "test/test_yaml_files/bad.yaml",
                "test/test_yaml_files/bad.yml",
                "test/test_yaml_files/good.yaml",
            ]
        result = main()
        self.assertEqual(result, True)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_false(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
                "test/test_yaml_files/good.yaml",
            ]
        result = main()
        self.assertEqual(result, False)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_no_changed_yaml_and_ignores_secret_directory(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "secret/yaml.yml", "some.yaml/but/not/a/yml.txt"
        ]
        result = main()
        self.assertEqual(result, False)
