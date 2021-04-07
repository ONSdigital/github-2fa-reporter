package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ONSdigital/github-2fa-reporter/pkg/github"
	"github.com/ONSdigital/github-2fa-reporter/pkg/slack"
)

func main() {
	enterprise := ""
	if enterprise = os.Getenv("GITHUB_ENTERPRISE_NAME"); len(enterprise) == 0 {
		log.Fatal("Missing GITHUB_ENTERPRISE_NAME environmental variable")
	}

	excludedLogins := ""
	if excludedLogins = os.Getenv("EXCLUDED_GITHUB_LOGINS"); len(excludedLogins) == 0 {
		log.Fatal("Missing EXCLUDED_GITHUB_LOGINS environmental variable")
	}

	token := ""
	if token = os.Getenv("GITHUB_TOKEN"); len(token) == 0 {
		log.Fatal("Missing GITHUB_TOKEN environmental variable")
	}

	slackAlertsChannel := ""
	if slackAlertsChannel = os.Getenv("SLACK_ALERTS_CHANNEL"); len(slackAlertsChannel) == 0 {
		log.Fatal("Missing SLACK_ALERTS_CHANNEL environment variable")
	}

	slackWebHookURL := ""
	if slackWebHookURL = os.Getenv("SLACK_WEBHOOK"); len(slackWebHookURL) == 0 {
		log.Fatal("Missing SLACK_WEBHOOK environment variable")
	}

	client := github.NewClient(token)
	outsideCollaborators, err := client.FetchOutsideCollaborators(enterprise)
	if err != nil {
		log.Fatalf("Failed to fetch outside collaborators: %v", err)
	}

	usersWithTwoFactorDisabled, err := client.FetchUsersWithTwoFactorDisabled(enterprise)
	if err != nil {
		log.Fatalf("Failed to fetch users with two-factor authentication disabled: %v", err)
	}

	exclusions := strings.Split(excludedLogins, "\n")
	text := "Good news: all of our human GitHub users have two-factor authentication enabled! :lock:"

	if len(usersWithTwoFactorDisabled) > 0 {
		var collaborators, members []github.User
		text = "The human GitHub users below do not have two-factor authentication enabled: :unlock:"

	USERS_LOOP:
		for _, user := range usersWithTwoFactorDisabled {
			outsideCollaborator := false

			for _, exclusion := range exclusions {
				if user.Login == exclusion {
					continue USERS_LOOP
				}
			}

			for _, id := range outsideCollaborators {
				if user.ID == id.ID {
					outsideCollaborator = true
					break
				}
			}

			if outsideCollaborator {
				collaborators = append(collaborators, user)
			} else {
				members = append(members, user)
			}
		}

		text = buildUsersText(members, "Enterprise Members", text)
		text = buildUsersText(collaborators, "Outside Collaborators", text)
	}

	postSlackMessage(text, slackAlertsChannel, slackWebHookURL)
}

func buildUsersText(users []github.User, heading, text string) string {
	if len(users) > 0 {
		text = fmt.Sprintf("%s\n\n*%s (%d):*\n", text, heading, len(users))

		for _, user := range users {
			text = fmt.Sprintf("%s\n%s", text, formatUser(user))
		}
	}

	return text
}

func formatUser(user github.User) string {
	displayName := fmt.Sprintf("<https://github.com/%s|%s>", user.Login, user.Login)

	if len(user.Name) > 0 && len(user.Email) > 0 {
		displayName = fmt.Sprintf("%s (%s, email: %s)", displayName, user.Name, user.Email)
	} else if len(user.Name) > 0 && len(user.Email) == 0 {
		displayName = fmt.Sprintf("%s (%s)", displayName, user.Name)
	} else if len(user.Name) == 0 && len(user.Email) > 0 {
		displayName = fmt.Sprintf("%s (email: %s)", displayName, user.Email)
	}

	return displayName
}

func postSlackMessage(text, channel, webHookURL string) {
	payload := slack.Payload{
		Text:      text,
		Username:  "GitHub 2FA Reporter",
		Channel:   channel,
		IconEmoji: ":github:",
	}

	err := slack.Send(webHookURL, payload)
	if err != nil {
		log.Fatalf("Failed to send Slack message: %v", err)
	}
}
