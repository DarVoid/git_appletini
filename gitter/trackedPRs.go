package gitter

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

type Resposta struct {
	data map[string]TrackedPullRequest
}

type TrackedPullRequest struct {
	label struct {
		pullRequests struct {
			edges []struct {
				node struct {
					id             string `yaml:"id"`
					title          string `yaml:"title"`
					url            string `yaml:"url"`
					baseRefName    string `yaml:"baseRefName"`
					headRefName    string `yaml:"headRefName"`
					reviewDecision string `yaml:"reviewDecision"`
					createdAt      string `yaml:"createdAt"`
					permalink      string `yaml:"permalink"`
					mergeable      string `yaml:"mergeable"`
					state          string `yaml:"state"`
					reviewRequests struct {
						totalCount int `yaml:"totalCount"`
					} `yaml:"reviewRequests"`
					reviews []struct {
						edges []struct {
							node struct {
								state  string `yaml:"state"`
								body   string `yaml:"body"`
								author struct {
									login string `yaml:"login"`
								} `yaml:"author"`
								comments []struct {
									edges []struct {
										node struct {
											body string `yaml:"body"`
										} `yaml:"node"`
									} `yaml:"edges"`
								} `yaml:"comments"`
							} `yaml:"node"`
						} `yaml:"edges"`
					} `yaml:"reviews"`
				} `yaml:"node"`
			} `yaml:"edges"`
		} `yaml:"pullRequests"`
	} `yaml:"label"`
}

func GetTrackedPullRequests(url string, query string, token string, ctx context.Context) (Resposta, error) {

	req := graphql.NewRequest(query)

	fmt.Println("Querying: ", url)
	client := graphql.NewClient(url)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	mapa := make(map[string]TrackedPullRequest, 0)
	resp := Resposta{
		data: mapa,
	}
	if err := client.Run(context.Background(), req, &mapa); err != nil {
		fmt.Println(err)
		return Resposta{}, err
	}
	fmt.Println("appletini: ", mapa)
	fmt.Println("appletini: ", mapa["appletini"])
	return resp, nil
}
