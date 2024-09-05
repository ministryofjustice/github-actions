# Slack - GitHub Secrets Scanning Alerts Integration

A GitHub Action to forward alerts from GitHub secret scanning and send them to Slack.

## Usage

```
      - uses: ministryofjustice/github-actions/slack-github-secret-scanning-integration@v18.1.4
        with:
          github-token: ${{ secrets.SECRET_SCANNING_GitHub_TOKEN }}
          slack-webhook-url: ${{ secrets.SLACK_WEBHOOK_URL }}
```

| Parameter                                                           | Description                                                      | Required                                                     | Default                                                     |
| ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- |
| frequency                                | Get secret scanning alerts that have occurred in this period prior to this action running | false | 24 hours |
| github-token                              | [GitHub token with access to secret scanning](https://docs.github.com/en/rest/secret-scanning/secret-scanning?apiVersion=2022-11-28#list-secret-scanning-alerts-for-a-repository)                 | true | NA |
| slack-webhook-url           | Incoming Slack webhook url for channel that you want to send alerts to                | true | NA |
