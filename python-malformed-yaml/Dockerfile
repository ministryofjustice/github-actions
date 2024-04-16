FROM python:3.12.0-alpine3.17
ADD . /app
WORKDIR /app

RUN pip install --target=/app PyYaml

COPY main.py /main.py

ENTRYPOINT [ "/app/main.py" ]


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