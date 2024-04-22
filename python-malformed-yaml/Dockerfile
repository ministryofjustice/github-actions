FROM python:3.12-slim

RUN apt-get update

WORKDIR /app

COPY python-malformed-yaml/github_service.py /app/github_service.py
COPY python-malformed-yaml/main.py /app/main.py
COPY python-malformed-yaml/requirements.txt /app/requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

ENTRYPOINT [ "python3", "/app/main.py" ]
