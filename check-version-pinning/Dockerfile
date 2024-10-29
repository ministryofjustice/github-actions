FROM python:3.12-slim

RUN apt-get update && apt-get upgrade -y && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY check_version_pinning.py /app/check_version_pinning.py

ENTRYPOINT ["python", "/app/check_version_pinning.py"]
CMD ["${{ inputs.workflow_directory }}", "${{ inputs.scan_mode }}"]
