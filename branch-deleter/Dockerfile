FROM ministryofjustice/cloud-platform-tools:1.4

RUN gem install octokit standardrb

COPY delete-branch-after-merge.rb /delete-branch-after-merge.rb
COPY github.rb /github.rb

ENTRYPOINT ["/delete-branch-after-merge.rb"]
