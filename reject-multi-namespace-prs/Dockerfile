FROM ministryofjustice/cloud-platform-tools:1.4

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit

COPY reject-multi-namespace-prs.rb /reject-multi-namespace-prs.rb
COPY github.rb /github.rb

ENTRYPOINT ["/reject-multi-namespace-prs.rb"]
