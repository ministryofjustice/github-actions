#!/bin/bash

WORKFLOW_DIRECTORY="${1:-.github/workflows}"
SCAN_MODE="${2:-full}"

python /app/check_version_pinning.py "$WORKFLOW_DIRECTORY" "$SCAN_MODE"
