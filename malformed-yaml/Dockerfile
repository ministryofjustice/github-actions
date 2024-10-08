FROM python:3.13-slim

RUN apt-get update && apt-get upgrade -y && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY github_pull_request.py main.py requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

ENTRYPOINT [ "python3", "/app/main.py" ]
