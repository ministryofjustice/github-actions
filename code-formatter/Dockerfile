FROM ruby:2.6-alpine

ENV \
  TERRAFORM_VERSION=0.13.6
ENV CFN_FORMATTER_VERSION=v1.1.2-1
# Install terraform
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip -d /usr/local/bin terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Install cloudformation-formatter
RUN apk add libc6-compat
RUN wget https://github.com/awslabs/aws-cloudformation-template-formatter/releases/download/${CFN_FORMATTER_VERSION}/cfn-format-${CFN_FORMATTER_VERSION}_linux-amd64.zip \
  && unzip -d /usr/local/bin cfn-format-${CFN_FORMATTER_VERSION}_linux-amd64.zip cfn-format-${CFN_FORMATTER_VERSION}_linux-amd64/cfn-format -j

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit standardrb

RUN apk add git

COPY format-code.rb /format-code.rb
COPY github.rb /github.rb
COPY code_formatter.rb /code_formatter.rb

ENTRYPOINT ["/format-code.rb"]
