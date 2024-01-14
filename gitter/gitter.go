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
					repository{
             branchProtectionRules(first:100){
        edges{
          node{
            allowsDeletions
            allowsForcePushes
          requiresApprovingReviews
            
          }
        }
      }
					  name
					  url
						owner{
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
            
					reviews(first:12) {
					  
					  totalCount
					  nodes{
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

type PrResponse struct {
	Viewer struct {
		PullRequests struct {
			Nodes []edge `json:"edges"`
		} `json:"pullRequests"`
	} `json:"viewer"`
}

// Gargabe stsart
type edge struct {
	Node pullRequest `json:"node"`
}

type pullRequest struct {
	Title       string `json:"title"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Number      int    `json:"number"`
	Permalink   string `json:"permalink"`
	ReviewCount struct {
		TotalCount int `json:"totalCount"`
	} `json:"reviewCount"`
	ReviewRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"reviewRequests"`
	ReviewDecision string `json:"reviewDecision"`
	Mergeable      string `json:"mergeable"`
}

// Garbage end
type PullRequest struct {
	Title       string `json:"title"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Number      int    `json:"number"`
	Permalink   string `json:"permalink"`
	ReviewCount int    `json:"reviewCount"`

	ReviewRequests int    `json:"reviewRequests"`
	ReviewDecision string `json:"reviewDecision"`

	Mergeable string `json:"mergeable"`
}

func (pr pullRequest) transform() PullRequest {

	return PullRequest{
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
