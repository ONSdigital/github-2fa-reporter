package github

import (
	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

// FetchOutsideCollaborators returns a list of outside collaborators within the enterprise.
func (c Client) FetchOutsideCollaborators(enterprise string) (outsideCollaborators []Identity, err error) {
	var endCursor *string // Using a pointer type allows this to be nil (an empty string isn't a valid cursor).

	req := graphql.NewRequest(`
		query ListOutsideCollaborators($slug: String!, $after: String) {
			enterprise(slug: $slug) {
				ownerInfo {
					outsideCollaborators(first: 50, after: $after) {
						pageInfo {
							startCursor
							endCursor
							hasNextPage
							hasPreviousPage
						}
						nodes {
							id
						}
					}
				}
			}
		}
	`)

	req.Var("slug", enterprise)

	hasNextPage := true
	var nodes []Identity

	for hasNextPage {
		res := &struct{ Enterprise Enterprise }{}
		req.Var("after", endCursor)

		if err := c.Run(req, &res); err != nil {
			return nil, errors.Wrap(err, "failed to fetch outside collaborators for enterprise")
		}

		nodes = append(nodes, res.Enterprise.OwnerInfo.OutsideCollaborators.Nodes...)
		endCursor = &res.Enterprise.OwnerInfo.OutsideCollaborators.PageInfo.EndCursor
		hasNextPage = res.Enterprise.OwnerInfo.OutsideCollaborators.PageInfo.HasNextPage
	}

	return nodes, nil
}
