FROM ministryofjustice/cloud-platform-tools:1.9

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit
RUN gem install colorize

COPY reject-escalated-privileges-yaml.rb /reject-escalated-privileges-yaml.rb
COPY github.rb /github.rb

RUN chmod +x /reject-escalated-privileges-yaml.rb

ENTRYPOINT ["/reject-escalated-privileges-yaml.rb"]
