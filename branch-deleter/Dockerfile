FROM ministryofjustice/cloud-platform-tools:2.1

RUN gem install octokit

COPY delete-branch-after-merge.rb /delete-branch-after-merge.rb
COPY github.rb /github.rb

ENTRYPOINT ["/delete-branch-after-merge.rb"]
