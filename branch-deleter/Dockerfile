FROM ministryofjustice/cloud-platform-tools:1.4

# Octokit depends on faraday, and an update to
# faraday breaks the current version of octokit
RUN gem install faraday --version 0.9
RUN gem install octokit

COPY delete-branch-after-merge.rb /delete-branch-after-merge.rb
COPY github.rb /github.rb

ENTRYPOINT ["/delete-branch-after-merge.rb"]
