FROM ministryofjustice/cloud-platform-tools:2.1

RUN gem install octokit

COPY reject-multi-namespace-prs.rb /reject-multi-namespace-prs.rb
COPY github.rb /github.rb

ENTRYPOINT ["/reject-multi-namespace-prs.rb"]
