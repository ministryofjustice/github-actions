FROM ministryofjustice/cloud-platform-tools:1.9

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit standardrb

COPY format-code.rb /format-code.rb
COPY github.rb /github.rb
COPY code_formatter.rb /code_formatter.rb

ENTRYPOINT ["/format-code.rb"]
