package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"git_applet/gitter"
	"os"
)

func makeHash(data []byte) string {
	cha := sha256.New()
	cha.Write(data)
	fmt.Printf("%x\n", string(cha.Sum(nil)))
	return fmt.Sprintf("%x", string(cha.Sum(nil)))
}

func hashPRs(prs []gitter.PullRequest) string {
	json_prs, err := json.Marshal(prs)
	ehp(err)
	return makeHash(json_prs)
}

func getPRs() []gitter.PullRequest {
	ctx := auth2()
	token := os.Getenv(Contexts[currentContext].Github.Token)
	prsLocal := gitter.PrResponse{}
	gqlApi := Contexts[currentContext].Github.GraphQL
	gitter.GetPullRequests(gqlApi, &prsLocal, token, ctx)
	fmt.Printf("prsLocal: %v\n", prsLocal)
	return prsLocal.Extract()
}

func polledPRs() {
	currentContext = config.DefaultContext
	Contexts = config.Contexts

	prs = getPRs()

	newHash := hashPRs(prs)

	if currentHash != newHash {
		currentHash = newHash
		syncPolledItems()
	}

}
func syncPolledItems() {
	clearPRItems()
	green, red := false, false
	for _, pr := range prs {
		pushPR(pr)

		fmt.Printf("pr.ReviewRequests: %v\n", pr.ReviewRequests)

		if pr.ReviewDecision == SHOW_GREEN_ON {
			green = true
		}

		if pr.ReviewDecision == SHOW_RED_ON {
			red = true
		}

		fmt.Printf("pr.ReviewDecision: %v\n", pr.ReviewDecision)

		if pr.ReviewCount > 0 {
			fmt.Printf("%v", "❗")
		}
	}

	if green && red {
		desk.SetSystemTrayIcon(resIconBoth)
		return
	}

	if red {
		desk.SetSystemTrayIcon(resIconReviewable)
		return
	}

	if green {
		desk.SetSystemTrayIcon(resIconMergeable)
		return
	}

	desk.SetSystemTrayIcon(resIconDefault)
}
