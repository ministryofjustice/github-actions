FROM python:3.12-slim

WORKDIR /app

COPY github_service.py /app/github_service.py
COPY reject_malformed_yaml.py /app/reject_malformed_yaml.py
COPY requirements.txt /app/requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

RUN chmod +x /app/reject_malformed_yaml.py

RUN python /app/reject_malformed_yaml.py
