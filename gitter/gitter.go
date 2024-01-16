package gitter

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func GetPullRequests(url string, data *PrResponse, token string, ctx context.Context) {

	req := graphql.NewRequest(
		`query fetchPRs {
			viewer {
			  pullRequests(
				orderBy: {field: CREATED_AT, direction: ASC}
				first: 100
				states: [OPEN]
			  ) {
				edges {
				  node {
					id
					repository {
					  branchProtectionRules(first: 100) {
						edges {
						  node {
							allowsDeletions
							allowsForcePushes
							requiresApprovingReviews
						  }
						}
					  }
					  name
					  url
					  owner {
						login
					  }
					}
					reviewDecision
					title
					baseRefName
					headRefName
					number
					permalink
					reviewRequests {
					  totalCount
					}
					reviews(first: 12) {
					  totalCount
					  nodes {
						state
					  }
					}
					mergeable
				  }
				}
			  }
			}
		  }`)

	client := graphql.NewClient(url)
	// TODO: do the same for organization

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	if err := client.Run(context.Background(), req, &data); err != nil {
		fmt.Println(err)
	}
}

func ApprovePullRequest(url string, token string, ctx context.Context, id string, body string) { //TODO: fix the graphql injection xD

	req := graphql.NewRequest(fmt.Sprintf(`mutation {
		addPullRequestReview(input: {
		  pullRequestId: "%s",
		  event: APPROVE,
		  body: "%s"
		}) {
		  pullRequestReview {
			id
			url
		  }
		}
	  }`, id, body))

	client := graphql.NewClient(url)
	// TODO: do the same for organization

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	if err := client.Run(context.Background(), req, nil); err != nil {
		fmt.Println(err)
	}
}

type PrResponse struct {
	Viewer struct {
		PullRequests struct {
			Nodes []edge `json:"edges"`
		} `json:"pullRequests"`
	} `json:"viewer"`
}

// Gargabe stsart
type edge struct {
	Node pullRequest `yaml:"node"`
}

type pullRequest struct {
	Id          string `yaml:"id"`
	Title       string `yaml:"title"`
	BaseRefName string `yaml:"baseRefName"`
	HeadRefName string `yaml:"headRefName"`
	Number      int    `yaml:"number"`
	Permalink   string `yaml:"permalink"`
	ReviewCount struct {
		TotalCount int `yaml:"totalCount"`
	} `yaml:"reviewCount"`
	ReviewRequests struct {
		TotalCount int `yaml:"totalCount"`
	} `yaml:"reviewRequests"`
	ReviewDecision string `yaml:"reviewDecision"`
	Mergeable      string `yaml:"mergeable"`
}

// Garbage end
type PullRequest struct {
	Title       string `yaml:"title"`
	BaseRefName string `yaml:"baseRefName"`
	HeadRefName string `yaml:"headRefName"`
	Number      int    `yaml:"number"`
	Permalink   string `yaml:"permalink"`
	ReviewCount int    `yaml:"reviewCount"`

	ReviewRequests int    `yaml:"reviewRequests"`
	ReviewDecision string `yaml:"reviewDecision"`
	Id             string `yaml:"id"`

	Mergeable  string         `yaml:"mergeable"`
	Remainder_ map[string]any `yaml:",inline"`
}

func (pr pullRequest) transform() PullRequest {

	return PullRequest{
		Id:             pr.Id,
		Title:          pr.Title,
		BaseRefName:    pr.BaseRefName,
		HeadRefName:    pr.HeadRefName,
		Number:         pr.Number,
		Permalink:      pr.Permalink,
		ReviewCount:    pr.ReviewCount.TotalCount,
		ReviewRequests: pr.ReviewRequests.TotalCount,
		ReviewDecision: pr.ReviewDecision,
		Mergeable:      pr.Mergeable,
	}
}
func (pr PrResponse) Extract() []PullRequest {
	prs := []PullRequest{}
	for _, val := range pr.Viewer.PullRequests.Nodes {
		prs = append(prs, val.Node.transform())
	}

	return prs
}
