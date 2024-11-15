package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"git_applet/gitter"
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
	token := getCurrentAccessToken()
	fmt.Println(token)
	prsLocal := gitter.PrResponse{}
	fmt.Printf("currentContext: %v\n", currentContext)
	gqlApi := getGraphQLApi()
	fmt.Printf("url: %v\n", gqlApi)

	gitter.GetPullRequests(gqlApi, &prsLocal, token, ctx)
	fmt.Printf("prsLocal: %v\n", prsLocal)
	return prsLocal.Extract()
}


func getCurrentAccessToken() string{
	return os.Getenv(Contexts[currentContext].Github.Token)
}

func getGraphQLApi() string{
	return Contexts[currentContext].Github.GraphQL
}

func polledPRs() {
	loadContext()
	for {
		time.Sleep(getPollDuration()) 
		prs = getPRs()
		fmt.Println(prs)
		newHash := hashPRs(prs)
		fmt.Println(newHash)
		if currentHash != newHash {
			currentHash = newHash
			syncPolledItems()
		}
	}

}

func getPollDuration() time.Duration{
	return time.Duration(Contexts[currentContext].Poll.Frequency * int(time.Second)) 
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

		// if pr.ReviewCount > 0 {
		// 	fmt.Printf("%v", "❗")
		// }
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
