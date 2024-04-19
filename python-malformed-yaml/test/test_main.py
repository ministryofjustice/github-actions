import pytest
import unittest
from unittest.mock import patch
from main import get_changed_yaml_files_from_pr, get_malformed_yaml_files, main
from github_service import GitHubService as github_service
# python_malformed_yaml_main = importlib.import_module(
#     "python-malformed-yaml.main")
# main = python_malformed_yaml_main.main

def test_get_changed_yaml_files_from_pr(mocker):
    # get_github_env_mock = mocker.patch("get_github_env")
    # get_github_env_mock.return_value = ("fake_github_token", "fake_repo", "123")
    # github_mock = mocker.patch("github")

    get_changed_files_from_pr_mock = mocker.patch("get_changed_files_from_pr")
    get_changed_files_from_pr_mock.return_value(["a.yml", "b.py", "c.yaml", "d.txt"])
    assert get_changed_yaml_files_from_pr() == ["a.yml", "c.yaml"]


def test_main_exception():
    with pytest.raises(Exception) as exc_info:
        main()
    expected_malformed_files = [
        "test_python_malformed_yaml/test_yaml_files/bad.yml",
        "test_python_malformed_yaml/test_yaml_files/bad.yaml"
    ]
    assert all(file in str(exc_info.value)
               for file in expected_malformed_files)
