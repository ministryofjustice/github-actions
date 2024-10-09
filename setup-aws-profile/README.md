# Setup AWS Profile Action

A GitHub Action to setup an aws profile.

## Usage

```
      - uses: ministryofjustice/github-actions/setup-aws-profile@v18.2.1
        with:
          role-arn: ${{ secrets.MY_AWS_ROLE_ARN }}
          profile-name: ${{ secrets.MY_PROFILE }}
```

| Parameter                                                           | Description                                                      | Required                                                     | Default                                                     |
| ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------- |
| role-arn                                | ARN of IAM role to create profile for | true | N/A |
| profile-name                              | Name of AWS profile                | true | N/A |
| aws-region           | AWS region                | false | eu-west-2 |
