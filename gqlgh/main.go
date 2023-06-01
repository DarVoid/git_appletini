package gqlgh

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"

	"golang.org/x/net/context"
)

type PrGQL struct {
	Data struct {
		Viewer struct {
			PullRequests struct {
				Nodes []edge `json:"edges"`
			} `json:"pullRequests"`
		} `json:"viewer"`
	} `json:"data"`
}

// Gargabe stsart
type edge struct {
	Node pullRequest `json:"node"`
}
type pullRequest struct {
	Title       string `json:"title"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Number      string `json:"number"`
	Permalink   string `json:"permalink"`
	ReviewCount struct {
		TotalCount int `json:"totalCount"`
	} `json:"reviewCount"`
	ReviewRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"reviewRequests"`
	ReviewDecision string `json:"reviewDecision"`
}

// Garbage end
type PullRequest struct {
	Title       string `json:"title"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Number      string `json:"number"`
	Permalink   string `json:"permalink"`
	ReviewCount int    `json:"reviewCount"`

	ReviewRequests int    `json:"reviewRequests"`
	ReviewDecision string `json:"reviewDecision"`
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
	}
}
func (pr PrGQL) Extract() []PullRequest {
	prs := []PullRequest{}
	for _, val := range pr.Data.Viewer.PullRequests.Nodes {
		prs = append(prs, val.Node.transform())
	}

	return prs
}

func GetPullRequests(url string, data *PrGQL, client *http.Client, ctx context.Context) {
	query := strings.NewReader(
		`{	"operationName":"fetchPRs",
			"variables":{},
			"query": 
			"query fetchPRs { viewer { pullRequests(orderBy: { field: CREATED_AT, direction: ASC}, first: 100 states: OPEN) { edges { node { title baseRefName headRefName number permalink reviewRequests { totalCount } reviews { totalCount } reviewDecision } } } } }"
			}`)

	resp, err := client.Post(url, "application/json", query)
	ehp(err)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, data)
}

func ehp(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
