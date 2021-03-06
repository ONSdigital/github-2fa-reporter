package github

import (
	"context"

	"github.com/ONSdigital/graphql"
)

type (

	// Client wraps a GraphQL client for communicating with the GitHub API.
	Client struct {
		token  string
		client *graphql.Client
	}

	// Enterprise represents the GitHub Enterprise.
	Enterprise struct {
		OwnerInfo struct {
			AffiliatedUsersWithTwoFactorDisabled struct {
				PageInfo PageInfo
				Nodes    []User
			}
			OutsideCollaborators struct {
				PageInfo PageInfo
				Nodes    []Identity
			}
		}
	}

	// PageInfo represents the pagination information returned from the query.
	PageInfo struct {
		StartCursor     string
		EndCursor       string
		HasPreviousPage bool
		HasNextPage     bool
	}

	// Identity represents the unique identity of an outside collaborator.
	Identity struct {
		ID string
	}

	// User represents the details of an affiliated user.
	User struct {
		ID    string
		Email string `json:"email,omitempty"`
		Login string `json:"login,omitempty"`
		Name  string `jsoh:"name,omitempty"`
	}
)

const endpoint = "https://api.github.com/graphql"

// NewClient instantiates a new GraphQL client.
func NewClient(token string) *Client {
	return &Client{
		token:  token,
		client: graphql.NewClient(endpoint),
	}
}

// Run wraps the underlying graphql.Run function, authomatically adding an authentication header and background context.
func (c Client) Run(request *graphql.Request, response interface{}) error {
	request.Header.Set("Authorization", "Bearer "+c.token)
	return c.client.Run(context.Background(), request, response)
}
