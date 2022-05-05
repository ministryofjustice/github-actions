FROM ministryofjustice/cloud-platform-tools:latest

RUN gem install octokit

COPY delete-branch-after-merge.rb /delete-branch-after-merge.rb
COPY github.rb /github.rb

ENTRYPOINT ["/delete-branch-after-merge.rb"]
