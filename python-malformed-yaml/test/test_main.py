import os
import logging
from unittest.mock import patch, MagicMock
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
    def test_get_changed_yaml_files_from_pr_no_yaml_system_exit(self, mock_github_service):
        mock_github_service.return_value.get_changed_files_from_pr.return_value = [
            "a.txt", "secret/non/yaml.csv", "some.yml/but/not/a/yml.txt"
        ]
        with self.assertRaises(SystemExit) as cm:
            get_changed_yaml_files_from_pr()
        self.assertEqual(cm.exception.code, 1)

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

    @patch.object(github_service, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_true(self, mock_github_service):
        mock_github_service.return_value.get_changed_files_from_pr.return_value = [
                "test/test_yaml_files/bad.yaml",
                "test/test_yaml_files/bad.yml",
                "test/test_yaml_files/good.yaml",
            ]
        result = main()
        self.assertEqual(result, True)

    @patch.object(github_service, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_false(self, mock_github_service):
        mock_github_service.return_value.get_changed_files_from_pr.return_value = [
                "test/test_yaml_files/good.yaml",
            ]
        result = main()
        self.assertEqual(result, False)
