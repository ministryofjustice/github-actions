FROM golang:1.16-buster

#Install Checkov
RUN apt-get update && apt-get install -y \
  python3.7 \
  python3-pip \
  git \
  jq \
  unzip \
  && rm -rf /var/lib/apt/lists/*
RUN pip3 install --upgrade pip && pip3 install --upgrade setuptools
RUN pip3 install checkov

#Install tflint
RUN curl https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
