package github

import (
	"sort"

	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

// FetchUsersWithTwoFactorDisabled returns a list of affiliated users with two-factor authentication disabled. The returned users are sorted by their login names.
func (c Client) FetchUsersWithTwoFactorDisabled(enterprise string) (users []User, err error) {
	var endCursor *string // Using a pointer type allows this to be nil (an empty string isn't a valid cursor).

	req := graphql.NewRequest(`
		query ListOutsideCollaborators($slug: String!, $after: String) {
			enterprise(slug: $slug) {
				ownerInfo {
					affiliatedUsersWithTwoFactorDisabled(first: 50, after: $after) {
						totalCount
						pageInfo {
							startCursor
							endCursor
							hasNextPage
							hasPreviousPage
						}
						nodes {
							id
							email
							login
							name
						}
					}
				}
			}
		}
	`)

	req.Var("slug", enterprise)

	hasNextPage := true
	var nodes []User

	for hasNextPage {
		res := &struct{ Enterprise Enterprise }{}
		req.Var("after", endCursor)

		if err := c.Run(req, &res); err != nil {
			return nil, errors.Wrap(err, "failed to fetch users with two-factor authentication disabled")
		}

		nodes = append(nodes, res.Enterprise.OwnerInfo.AffiliatedUsersWithTwoFactorDisabled.Nodes...)
		endCursor = &res.Enterprise.OwnerInfo.AffiliatedUsersWithTwoFactorDisabled.PageInfo.EndCursor
		hasNextPage = res.Enterprise.OwnerInfo.AffiliatedUsersWithTwoFactorDisabled.PageInfo.HasNextPage
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Login < nodes[j].Login
	})

	return nodes, nil
}
