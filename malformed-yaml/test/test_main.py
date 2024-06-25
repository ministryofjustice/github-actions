import os
from unittest.mock import patch
import unittest
from main import (
    get_changed_yaml_files_from_pr,
    get_extant_files,
    get_malformed_yaml_files_and_errors,
    main
    )
from github_pull_request import GitHubPullRequest as github_pull_request


class TestMain(unittest.TestCase):

    def test_get_extant_files(self):
        files = [
                "random/non_existant_file.txt"
                "test/test_yaml_files/deleted_file.yaml",
                "test/test_yaml_files/bad.yml",
                "test/test_yaml_files/secret/good.yaml"
            ]
        result = get_extant_files(files=files)
        self.assertIn("test/test_yaml_files/bad.yml", "\n".join(result))
        self.assertIn("test/test_yaml_files/secret/good.yaml", "\n".join(result))
        self.assertNotIn("test/test_yaml_files/deleted_file.yaml", "\n".join(result))
        self.assertNotIn("random/non_existant_file.txt", "\n".join(result))

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
                "test/test_yaml_files/good-multi-doc.yml",
                "test/test_yaml_files/bad-multi-doc.yaml",
            ]
        )
        self.assertIn("test/test_yaml_files/bad.yaml", "\n".join(result))
        self.assertIn("test/test_yaml_files/bad.yml", "\n".join(result))
        self.assertNotIn("test/test_yaml_files/good.yaml", "\n".join(result))
        self.assertNotIn("test/test_yaml_files/good-multi-doc.yml", "\n".join(result))
        self.assertIn("test/test_yaml_files/bad-multi-doc.yaml", "\n".join(result))

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_true(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
                "test/skip/deleted_file.yaml",
                "test/test_yaml_files/bad.yaml",
                "test/test_yaml_files/bad.yml",
                "test/test_yaml_files/good.yaml",
                "test/test_yaml_files/good-multi-doc.yml",
                "test/test_yaml_files/bad-multi-doc.yaml",
                "test/test_yaml_files/secret/bad.yml"
            ]
        result = main()
        self.assertEqual(result, True)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_malformed_yaml_false(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
                "test/skip/deleted_file.yaml",
                "test/test_yaml_files/good.yaml",
                "test/test_yaml_files/good-multi-doc.yml",
                "test/test_yaml_files/secret/bad.yml"
            ]
        result = main()
        self.assertEqual(result, False)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_no_yaml_changes_and_ignores_secret_directory(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "test/test_yaml_files/secret/good.yaml",
            "test/test_yaml_files/secret/bad.yml"
        ]
        result = main()
        self.assertEqual(result, False)

    @patch.object(github_pull_request, "__new__")
    @patch.dict(os.environ, {"GITHUB_TOKEN": "token", "PR_NUMBER": "123", "REPOSITORY_NAME": "repo_name"})    
    def test_main_no_yaml_changes_and_ignores_deleted_files(self, mock_github_pull_request):
        mock_github_pull_request.return_value.get_changed_files_from_pr.return_value = [
            "test/skip/deleted_file.yaml",
        ]
        result = main()
        self.assertEqual(result, False)
