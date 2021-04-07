# GitHub 2FA Reporter
This repository contains a [Go](https://golang.org/) application that consumes the [GitHub GraphQL API](https://docs.github.com/en/graphql) and posts a Slack message containing details of GitHub users with two-factor authentication disabled. Both enterprise members and outside collaborators are included.

## Building
Use `make` to compile binaries for macOS and Linux.

## Running
### Environment Variables
The environment variables below are required:

```
EXCLUDED_GITHUB_LOGINS # List of GitHub logins with two-factor authentication disabled to exclude from reporting
GITHUB_ENTERPRISE_NAME # Name of the GitHub Enterprise
GITHUB_TOKEN           # GitHub personal access token
SLACK_ALERTS_CHANNEL   # Name of the Slack channel to post alerts to
SLACK_WEBHOOK          # Used for accessing the Slack Incoming Webhooks API
```

### Token Scopes
The GitHub personal access token for using this application requires the following scopes:

- `admin:enterprise`
- `read:user`
- `user:email`

## Copyright
Copyright (C) 2021 Crown Copyright (Office for National Statistics)