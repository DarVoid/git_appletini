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
	fmt.Println("YO")
	decisions := make(map[string]string)
	decisions["APPROVED"] = "APPROVED ✔️"
	decisions["CHANGES_REQUESTED"] = "RIP, you tried💩"
	decisions[""] = "on Hold..."

	// prBox.RemoveChildren()
	green, red := false, false
	for _, pr := range prs {
		fmt.Printf("pr.ReviewRequests: %v\n", pr.ReviewRequests)
		if pr.ReviewDecision == "APPROVED" {
			green = true
		}
		if pr.ReviewDecision == "" {
			red = true
		}
		fmt.Printf("pr.ReviewDecision: %v\n", pr.ReviewDecision)
		// status, ok := decisions[pr.ReviewDecision]

		// if !ok {
		// 	status = ""
		// }
		fmt.Println("YO2")

		if pr.ReviewCount > 0 {
			fmt.Printf("%v", "❗")
		}
		fmt.Println("YO3")

		// val := fmt.Sprintf("%v => %v\nStatus:%v\n\n -- %v", pr.HeadRefName, pr.BaseRefName, status, pr.Title)
		// item := prBox.AddSubMenuItem(val, val)
		// go handleLink(item, pr.Permalink) // TODO: uncomment
	}
	if green && red {
		fmt.Println("YOsad")

		desk.SetSystemTrayIcon(resIconBoth)
		return
	}
	if red {
		fmt.Println("YOsad2")
		desk.SetSystemTrayIcon(resIconReviewable)
		return
	}
	if green {
		fmt.Println("YOsad3")
		desk.SetSystemTrayIcon(resIconMergeable)
		return
	}
	fmt.Println("YO4")

	desk.SetSystemTrayIcon(resIconDefault)
}
