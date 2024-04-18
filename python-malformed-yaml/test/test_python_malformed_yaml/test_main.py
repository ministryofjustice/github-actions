import pytest
import importlib
python_malformed_yaml_main = importlib.import_module(
    "python-malformed-yaml.main")
main = python_malformed_yaml_main.main


def test_main_exception():
    with pytest.raises(Exception) as exc_info:
        main()
    expected_malformed_files = [
        "test_python_malformed_yaml/test_yaml_files/bad.yml",
        "test_python_malformed_yaml/test_yaml_files/bad.yaml"
    ]
    assert all(file in str(exc_info.value)
               for file in expected_malformed_files)
