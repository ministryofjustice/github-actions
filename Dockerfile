FROM ministryofjustice/cloud-platform-tools:1.4

RUN gem install octokit standardrb

COPY format-code.rb /format-code.rb

ENTRYPOINT ["/format-code.rb"]
