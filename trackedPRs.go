package main

import (
	"fmt"
	"git_applet/gitter"
	"time"
)

func polledTrackedPRs() {
	loadContext()
	for {
		time.Sleep(getPollDuration())
		prs, err := getTrackedPRs()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("PRS:", prs)
		// newHash := hashPRs(prs)
		// fmt.Println(newHash)
		// if currentTrackedHash != newHash {
		// 	currentTrackedHash = newHash
		// 	syncPolledTrackedItems()
		// }
	}

}
func getTrackedPRs() (any, error) {
	ctx := auth2()
	token := getCurrentAccessToken()
	fmt.Println(token) //maybe this needs changing
	fmt.Printf("currentContext: %v\n", currentContext)
	gqlApi := getGraphQLApi()
	fmt.Printf("url: %v\n", gqlApi)

	prs, err := gitter.GetTrackedPullRequests(gqlApi, SavedPRQuerry, token, ctx)
	if err != nil {
		fmt.Println("deu merda")
	}
	return prs, nil
	// return prsLocal.Extract()
}

// visual stuff

func syncPolledTrackedItems() {
	clearTrackedPRItems()
	green, red := false, false
	for _, pr := range prs {
		pushTrackedPR(pr)

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
