FROM ministryofjustice/cloud-platform-tools:2.7.0

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 2.7.4
RUN gem install octokit --version 6.1.0

COPY reject-multi-namespace-prs.rb /reject-multi-namespace-prs.rb
COPY github.rb /github.rb

ENTRYPOINT ["/reject-multi-namespace-prs.rb"]
