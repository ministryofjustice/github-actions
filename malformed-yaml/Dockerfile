FROM python:3.12-slim

WORKDIR /app
COPY . .
RUN pip install --no-cache-dir -r requirements.txt

ENTRYPOINT ["python", "reject_malformed_yaml.py"]