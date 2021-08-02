FROM golang:1.16-buster

# Install Terraform
RUN apt update && \
  apt install -y software-properties-common && \
  apt-add-repository "deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main" && \
  curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - && \
  apt update && \
  apt install -y terraform

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
