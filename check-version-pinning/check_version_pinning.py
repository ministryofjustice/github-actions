import json
import os
import sys

import yaml


def find_workflow_files(workflow_directory):
    for root, _, files in os.walk(workflow_directory):
        for file in files:
            if file.endswith(".yml") or file.endswith(".yaml"):
                yield os.path.join(root, file)


def find_changed_files_in_pr(workflow_directory):
    event_path = os.getenv("GITHUB_EVENT_PATH")
    if not event_path:
        print("Error: GITHUB_EVENT_PATH is not set.")
        sys.exit(1)

    with open(event_path, "r") as f:
        event_data = json.load(f)

    changed_files = [
        file["filename"]
        for file in event_data.get("pull_request", {}).get("files", [])
        if file["filename"].startswith(workflow_directory)
        and (file["filename"].endswith(".yml") or file["filename"].endswith(".yaml"))
    ]

    return changed_files


def parse_yaml_file(file_path):
    with open(file_path, "r", encoding="utf-8") as f:
        try:
            return yaml.safe_load(f)
        except yaml.YAMLError as e:
            print(f"Error parsing {file_path}: {e}")
            return None


def check_uses_field_in_workflow(workflows, file_path):
    results = []
    if workflows:
        for job in workflows.get("jobs", {}).values():
            for step in job.get("steps", []):
                uses = step.get("uses", "")
                if "@v" in uses:
                    results.append(f"{file_path}: {uses}")
    return results


def check_version_pinning(workflow_directory=".github/workflows", scan_mode="full"):
    all_results = []

    if scan_mode == "full":
        files_to_check = find_workflow_files(workflow_directory)
    elif scan_mode == "pr_changes":
        files_to_check = find_changed_files_in_pr(workflow_directory)
    else:
        print("Error: Invalid scan mode. Choose 'full' or 'pr_changes'.")
        sys.exit(1)

    for file_path in files_to_check:
        workflows = parse_yaml_file(file_path)
        if workflows:
            results = check_uses_field_in_workflow(workflows, file_path)
            all_results.extend(results)

    if all_results:
        print(
            "The following GitHub Actions are using version pinning rather than SHA hash pinning:\n"
        )
        for result in all_results:
            print(f"  - {result}")

        print(
            "\nPlease see the following documentation for more information:\n"
            "https://tinyurl.com/3sev9etr"
        )
        sys.exit(1)
    else:
        print("No workflows found with pinned versions (@v).")


if __name__ == "__main__":
    workflow_directory = sys.argv[1] if len(sys.argv) > 1 else ".github/workflows"
    scan_mode = sys.argv[2] if len(sys.argv) > 2 else "full"
    check_version_pinning(workflow_directory, scan_mode)
