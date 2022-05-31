FROM ruby:2.6-alpine

ENV CONFTEST_VERSION=0.32.0

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit --version 4.21.0
RUN gem install standardrb

# Install conftest binary
RUN wget https://github.com/open-policy-agent/conftest/releases/download/v${CONFTEST_VERSION}/conftest_${CONFTEST_VERSION}_Linux_x86_64.tar.gz \
  && tar xzf conftest_${CONFTEST_VERSION}_Linux_x86_64.tar.gz \
  && mv conftest /usr/local/bin \
  && rm conftest_${CONFTEST_VERSION}_Linux_x86_64.tar.gz

# Install opa binary
RUN wget https://www.openpolicyagent.org/downloads/latest/opa_linux_amd64 \
  && mv opa_linux_amd64 /usr/local/bin/opa \
  && chmod +x /usr/local/bin/opa

COPY run-conftest.rb /run-conftest.rb
COPY github.rb /github.rb

ENTRYPOINT ["/run-conftest.rb"]