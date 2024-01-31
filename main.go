package main

import (
	_ "embed"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"git_applet/actions"
	"git_applet/gitter"
)

func main() {
	loadConfig()
	loadContext()
	a := createApp()
	setupItems()
	go polledPRs()
	a.Run()
}

func loadContext() {
	currentContext = config.DefaultContext
	Contexts = config.Contexts
}

func createApp() fyne.App {
	a := app.NewWithID("git_appletini")
	var ok bool
	desk, ok = a.(desktop.App)
	if !ok {
		panic("could not create desktop app")
	}

	mprincipal = fyne.NewMenu("Git Applet")
	desk.SetSystemTrayMenu(mprincipal)

	a.Lifecycle().SetOnStarted(func() {
		desk.SetSystemTrayIcon(resIconDefault)
	})
	return a
}

func pushPR(pr gitter.PullRequest) {

	approve_status, _ := decision_messages[pr.ReviewDecision]
	merge_status, _ := merge_messages[pr.Mergeable]

	title := fmt.Sprintf("🔷 (#%d) %s\n[%s] ↦ [%s]\n%s\n%s", pr.Number, pr.Title, pr.HeadRefName, pr.BaseRefName, approve_status, merge_status)

	fmt.Println("author: ", pr.Author)
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
			token := os.Getenv(Contexts[currentContext].Github.Token)
			gqlApi := Contexts[currentContext].Github.GraphQL
			gitter.ApprovePullRequest(gqlApi, token, ctx, pr.Id, "LGTM! 🚀")
		},
	})
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
		fyne.NewMenuItem("Change icon", func() {
			fmt.Println("Clicked")
			desk.SetSystemTrayIcon(resIconReviewable)
		}), // TODO: remove when app complete (showcase only)
		fyne.NewMenuItem("Delete self", func() {
			fmt.Println("Clicked")
			mprincipal.Items = mprincipal.Items[:1] // how to delete stuff
			Refresh()
		}), // TODO: remove when app complete (showcase only)
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
