FROM ministryofjustice/cloud-platform-tools:1.9

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit

COPY reject-malformed-yaml.rb /reject-malformed-yaml.rb
COPY github.rb /github.rb

ENTRYPOINT ["/reject-malformed-yaml.rb"]
