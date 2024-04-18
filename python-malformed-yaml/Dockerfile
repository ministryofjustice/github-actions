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


# FROM python:3.12-slim

# # Setup env
# ENV LANG C.UTF-8
# ENV LC_ALL C.UTF-8
# ENV PYTHONDONTWRITEBYTECODE 1
# ENV PYTHONFAULTHANDLER 1

# FROM base AS python-deps

# # Install pipenv and compilation dependencies
# RUN pip install pipenv

# # Install python dependencies in /.venv
# COPY Pipfile .
# COPY Pipfile.lock .
# RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy

# FROM base AS runtime

# # Copy virtual env from python-deps stage
# COPY --from=python-deps /.venv /.venv
# ENV PATH="/.venv/bin:$PATH"

# # Create and switch to a new user
# RUN useradd --create-home appuser
# WORKDIR /app
# USER appuser

# # Install files into container
# COPY github_service.py /app/github_service.py
# COPY main.py /app/main.py

# # Run 
# RUN python main.py


FROM python:3.12-slim

RUN apt-get update

WORKDIR /app

COPY Pipfile /app/Pipfile
COPY Pipfile.lock /app/Pipfile.lock
COPY github_service.py /app/github_service.py
COPY main.py /app/main.py

RUN pip install pipenv
RUN pipenv install

RUN pipenv run python /app/main.py
