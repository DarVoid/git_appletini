package main

import (
	_ "embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"git_applet/actions"
	"git_applet/gitter"
)

func main() {
	loadConfig()
	fmt.Println(Config)
	loadContext()
	fmt.Println(Contexts)
	a := createApp()
	setupItems()
	go polledPRs()
	a.Run()
}

func loadContext() {
	currentContext = Config.DefaultContext
	Contexts = Config.Contexts
}

func createApp() fyne.App {
	// create app
	a := app.NewWithID("git_appletini")
	var ok bool
	// get desktop app
	desk, ok = a.(desktop.App)
	if !ok {
		panic("could not create desktop app")
	}
	// create main menu
	mprincipal = fyne.NewMenu("Git Applet")
	// set main menu onto the desktop app created
	desk.SetSystemTrayMenu(mprincipal)

	a.Lifecycle().SetOnStarted(func() {
		// set default icon on main menu (github white simple)
		desk.SetSystemTrayIcon(resIconDefault)
	})
	return a
}

func pushPR(pr gitter.PullRequest) {

	approveStatus, _ := decision_messages[pr.ReviewDecision]

	mergeStatus := checkIfMergeable(pr.Mergeable, pr.ReviewDecision)

	title := fmt.Sprintf("🔷 (#%d) %s\n[%s] ↦ [%s]\n%s\n%s", pr.Number, pr.Title, pr.HeadRefName, pr.BaseRefName, approveStatus, mergeStatus)

	pushPRItem(title, map[string]func(){
		"Open in browser": func() {
			actions.OpenLink(pr.Permalink, Contexts[currentContext].ChromeProfile)
		},
		"Close PR": func() {
			// TODO: Close PR
			fmt.Printf("%v\n", pr.Remainder_)
		},
		"Auto-reply with \"LGTM\"": func() {
			//! Please don't
			ctx := auth2()
			token := getCurrentAccessToken()
			gqlApi := getGraphQLApi()
			gitter.ApprovePullRequest(gqlApi, token, ctx, pr.Id, "LGTM! 🚀")
		},
	})
}

func checkIfMergeable(mergeableStatus string, approveStatus string) string {
	fmt.Println(mergeableStatus, approveStatus)
	// default message to be shown should be under this key
	decision := "NO_BUENO"
	// combination of different statuses set decision
	if mergeableStatus == "MERGEABLE" && approveStatus == "REVIEW_REQUIRED" {
		decision = "REQUIRES_REVIEW"
	}
	return merge_messages[decision]
}

func pushPRItem(title string, actions map[string]func()) {
	prItem := fyne.NewMenuItem(title, func() {})

	prItem.ChildMenu = fyne.NewMenu("Actions")
	for name, action := range actions {
		prItem.ChildMenu.Items = append(
			prItem.ChildMenu.Items,
			fyne.NewMenuItem(name, action),
		)
	}

	prBox.ChildMenu.Items = append(
		prBox.ChildMenu.Items,
		prItem,
		fyne.NewMenuItemSeparator(),
	)
	Refresh()
}

func clearPRItems() {
	prBox.ChildMenu.Items = []*fyne.MenuItem{}
	Refresh()
}

func setupItems() {
	setupPRBox()
	setupContextSelector()

	mprincipal.Items = []*fyne.MenuItem{
		// fyne.NewMenuItem("Change icon", func() {
		// 	fmt.Println("Clicked")
		// 	desk.SetSystemTrayIcon(resIconReviewable)
		// }),
		// fyne.NewMenuItem("Delete self", func() {
		// 	fmt.Println("Clicked")
		// 	mprincipal.Items = mprincipal.Items[:1] // how to delete stuff
		// 	Refresh()
		// }),
		fyne.NewMenuItemSeparator(),
		prBox,
		fyne.NewMenuItemSeparator(),
		contextSelector,
	}
}

func setupPRBox() {
	prBox = fyne.NewMenuItem("📑 Pull Requests", func() {})
	prBox.ChildMenu = fyne.NewMenu("Pull Requests")
}

func makeContextLabel() string {
	return fmt.Sprintf("👥 Context: %s", Contexts[currentContext].Title)
}

func setupContextSelector() {
	contextSelector = fyne.NewMenuItem(makeContextLabel(), func() {})
	contextSelector.ChildMenu = fyne.NewMenu("Context")
	for key, context := range Contexts {
		key := key
		contextSelector.ChildMenu.Items = append(contextSelector.ChildMenu.Items, fyne.NewMenuItem(context.Title, func() {
			currentContext = key
			contextSelector.Label = makeContextLabel()
			Refresh()
		}))
	}
}

func ehp(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
}
func Refresh() {
	mprincipal.Refresh()
}
