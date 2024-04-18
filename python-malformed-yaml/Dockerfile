# FROM python:3-slim AS builder
# ADD . /app
# WORKDIR /app

# # We are installing a dependency here directly into our app source dir
# RUN pip install --target=/app PyYaml

# # A distroless container image with Python and some basics like SSL certificates
# # https://github.com/GoogleContainerTools/distroless
# FROM gcr.io/distroless/python3-debian10
# COPY --from=builder /app /app
# WORKDIR /app
# ENV PYTHONPATH /app
# CMD ["/app/main.py"]


FROM python:3.12-slim

WORKDIR /app

COPY github_service.py /app/github_service.py
COPY main.py /app/main.py
COPY Pipfile /app/Pipfile
COPY Pipfile.lock /app/Pipfile.lock

RUN python3 -m pip install pipenv
RUN pipenv install
RUN chmod +x /app/main.py
RUN pipenv run python3 /app/main.py