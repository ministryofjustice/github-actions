# Send GitHub advisories to Slack

A Github Action to produce a Slack message containing GitHub advisories to be sent to a 
Slack channel (e.g. using [slackapi/slack-github-action](https://github.com/slackapi/slack-github-action)).

When configured as below, it will run regularly and post advisories to the provided channel.

## Usage

  1. [Set up or choose a Slack App](https://api.slack.com/apps?new_app=1) 
  that will post the Slack alerts into a specified channel. 
  Add the `chat:write` scope to the bot token (within `OAuth & Permissions`) 
  and then submit the Slack App for review.

  2. Add the Slack App 'Bot User OAuth Token' (within `OAuth & Permissions`) to your repo's
  environment secrets (available to actions) - below it is referenced as `SLACK_BOT_TOKEN`

  3. Create a file in your repo called `.github/workflows/format-github-advisories-for-slack.yml` 
  with the following contents:

  ```
  name: Send [Application] security advisory alerts to Slack

  on:
    schedule:
      - cron: '0 0 * * *'  # Runs every day at midnight
    workflow_dispatch:

  jobs:
    check-security-advisories:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4

        - name: Format advisories for target repo
          uses: ministryofjustice/github-actions/format-github-advisories-for-slack@v18
          with:
              target-repo-owner: "github-owner"
              target-repo: "target-repo-name"
              version-file-path: "helm_deploy/values.yaml"
              version-key: "global.application.version"
          env:
            GH_TOKEN: ${{ github.token }}
    
        - name: Send advisories to Slack
          id: slack
          uses: slackapi/slack-github-action@v1.26.0
          with:
            channel-id: "CXXXXXXXXXX" # Channel > Right click > View Channel Details > Scroll to bottom > Channel ID
            payload-file-path: "./filtered_advisories.json"
          env:
            SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
  ```

  `GH_TOKEN` is provided automatically by github actions. 
  You do not need to do anything extra to make it available.
