FROM ruby:2.6-alpine

ENV \
  TERRAFORM_VERSION=0.12.17

# Install terraform
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip -d /usr/local/bin terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit standardrb

COPY format-code.rb /format-code.rb
COPY github.rb /github.rb
COPY code_formatter.rb /code_formatter.rb

ENTRYPOINT ["/format-code.rb"]
