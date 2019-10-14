FROM ministryofjustice/cloud-platform-tools:1.4

RUN gem install octokit standardrb

WORKDIR /app

COPY bin /app/.github-action-bin

ENTRYPOINT "./.github-action-bin/format-code.rb"
