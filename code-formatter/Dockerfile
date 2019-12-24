FROM ministryofjustice/cloud-platform-tools:1.5

RUN gem install octokit standardrb

COPY format-code.rb /format-code.rb
COPY github.rb /github.rb
COPY code_formatter.rb /code_formatter.rb

ENTRYPOINT ["/format-code.rb"]
