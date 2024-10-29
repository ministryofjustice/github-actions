import unittest
from unittest.mock import mock_open, patch

from bin.check_version_pinning import check_version_pinning


class TestCheckVersionPinning(unittest.TestCase):

    @patch("os.walk")
    @patch("builtins.open", new_callable=mock_open)
    @patch("yaml.safe_load")
    def test_no_yaml_files(self, mock_yaml_load, mock_open_file, mock_os_walk):
        # Simulate os.walk returning no .yml or .yaml files
        _ = mock_open_file
        mock_os_walk.return_value = [(".github/workflows", [], [])]
        mock_yaml_load.return_value = None

        with patch("builtins.print") as mock_print:
            check_version_pinning()
            mock_print.assert_called_once_with(
                "No workflows found with pinned versions (@v)."
            )

    @patch("os.walk")
    @patch("builtins.open", new_callable=mock_open)
    @patch("yaml.safe_load")
    def test_yaml_file_without_uses(self, mock_yaml_load, mock_open_file, mock_os_walk):
        _ = mock_open_file
        mock_os_walk.return_value = [(".github/workflows", [], ["workflow.yml"])]
        mock_yaml_load.return_value = {
            "jobs": {"build": {"steps": [{"name": "Checkout code"}]}}
        }

        with patch("builtins.print") as mock_print:
            check_version_pinning()
            mock_print.assert_called_once_with(
                "No workflows found with pinned versions (@v)."
            )

    @patch("os.walk")
    @patch("builtins.open", new_callable=mock_open)
    @patch("yaml.safe_load")
    def test_workflow_with_pinned_version(
        self, mock_yaml_load, mock_open_file, mock_os_walk
    ):
        # Simulate a workflow file with a pinned version (@v)
        _ = mock_open_file
        mock_os_walk.return_value = [(".github/workflows", [], ["workflow.yml"])]
        mock_yaml_load.return_value = {
            "jobs": {"build": {"steps": [{"uses": "some-org/some-action@v1.0.0"}]}}
        }

        with patch("builtins.print") as mock_print, self.assertRaises(SystemExit) as cm:
            check_version_pinning()
            mock_print.assert_any_call("Found workflows with pinned versions (@v):")
            mock_print.assert_any_call(
                ".github/workflows/workflow.yml: some-org/some-action@v1.0.0"
            )
            self.assertEqual(cm.exception.code, 1)

    @patch("os.walk")
    @patch("builtins.open", new_callable=mock_open)
    @patch("yaml.safe_load")
    def test_workflow_ignoring_actions(
        self, mock_yaml_load, mock_open_file, mock_os_walk
    ):
        _ = mock_open_file
        # Simulate a workflow file with an action to be ignored
        mock_os_walk.return_value = [(".github/workflows", [], ["workflow.yml"])]
        mock_yaml_load.return_value = {
            "jobs": {
                "build": {
                    "steps": [
                        {"uses": "actions/setup-python@v2"},
                        {"uses": "ministryofjustice/some-action@v1.0.0"},
                    ]
                }
            }
        }

        with patch("builtins.print") as mock_print:
            check_version_pinning()
            mock_print.assert_called_once_with(
                "No workflows found with pinned versions (@v)."
            )

    @patch("os.walk")
    @patch("builtins.open", new_callable=mock_open)
    @patch("yaml.safe_load")
    def test_workflow_with_mixed_versions(
        self, mock_yaml_load, mock_open_file, mock_os_walk
    ):
        _ = mock_open_file
        # Simulate a workflow with both ignored and non-ignored actions
        mock_os_walk.return_value = [(".github/workflows", [], ["workflow.yml"])]
        mock_yaml_load.return_value = {
            "jobs": {
                "build": {
                    "steps": [
                        {"uses": "actions/setup-python@v2"},
                        {"uses": "some-org/some-action@v1.0.0"},
                        {"uses": "ministryofjustice/some-action@v1.0.0"},
                    ]
                }
            }
        }

        with patch("builtins.print") as mock_print, self.assertRaises(SystemExit) as cm:
            check_version_pinning()
            mock_print.assert_any_call("Found workflows with pinned versions (@v):")
            mock_print.assert_any_call(
                ".github/workflows/workflow.yml: some-org/some-action@v1.0.0"
            )
            self.assertEqual(cm.exception.code, 1)


if __name__ == "__main__":
    unittest.main()
