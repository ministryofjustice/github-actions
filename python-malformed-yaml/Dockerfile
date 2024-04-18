FROM python:3.12-slim

RUN apt-get update

WORKDIR /app

COPY Pipfile /app/Pipfile
COPY Pipfile.lock /app/Pipfile.lock
COPY github_service.py /app/github_service.py
COPY main.py /app/main.py

RUN pip install pipenv
RUN pipenv install

ENTRYPOINT [ "pipenv", "run", "python3", "/app/main.py" ]
