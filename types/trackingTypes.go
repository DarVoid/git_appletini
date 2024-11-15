package types

import (
	"fmt"
	"strings"
)

type TrackingConfig struct {
	RepoList        RepoItemSet        `json:"repoList"`
	LabeledRepoList LabeledRepoItemSet `json:"labeledRepoList"`
	ReviewAmmount   uint               `json:"reviewAmmount"`
	PrAmmount       uint               `json:"prAmmount"`
	CommentsAmmount uint               `json:"commentsAmmount"`
}

type RepoItem struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Identifier string `json:"identifier"`
}

type LabeledRepoItem struct {
	RepoItem
	Label string `json:"label"`
}
type LabeledRepoItemSet []LabeledRepoItem
type RepoItemSet []RepoItem

func (tr TrackingConfig) String() string {
	return fmt.Sprintf(`{
	LabeledRepoList: %v,
	RepoList: %v,
	PrAmmount: %v,
	ReviewAmmount: %v,
	CommentsAmmount: %v
	}`,
		tr.LabeledRepoList,
		tr.RepoList,
		tr.PrAmmount,
		tr.ReviewAmmount,
		tr.CommentsAmmount)
}
func (repo LabeledRepoItem) String() string {
	return fmt.Sprintf(`{
		Owner: %v,
		Name: %v,
		Identifier: %v,
		Label: %v
	}`,
		repo.Owner,
		repo.Name,
		repo.Identifier,
		repo.Label)
}
func (repo RepoItem) String() string {
	return fmt.Sprintf(`{
		Owner: %v,
		Name: %v,
		Identifier: %v
	}`,
		repo.Owner,
		repo.Name,
		repo.Identifier)
}
func (repo RepoItemSet) String() string {
	elems := []string{}
	for _, a := range repo {
		elems = append(elems, fmt.Sprint(a))
	}
	return fmt.Sprint("[", strings.Join(elems, ",\n\t"), "]")
}

func (repo LabeledRepoItemSet) String() string {
	elems := []string{}
	for _, a := range repo {
		elems = append(elems, fmt.Sprint(a))
	}
	return fmt.Sprint("[", strings.Join(elems, ",\n\t"), "]")
}
