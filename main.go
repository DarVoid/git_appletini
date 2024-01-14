package main

import (
	_ "embed"
	"fmt"
	"git_applet/actions"
	"git_applet/gitter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	a := createApp()
	addItems()
	go polledPRs()
	a.Run()
}

func createApp() fyne.App {
	a := app.NewWithID("git_appletini")
	loadConfig()
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

func clearPRItems() {
	prBox.ChildMenu.Items = []*fyne.MenuItem{}
	mprincipal.Refresh()
}

func pushPRItem(title string, action func()) {
	prBox.ChildMenu.Items = append(
		prBox.ChildMenu.Items,
		fyne.NewMenuItem(title, action),
		fyne.NewMenuItemSeparator(),
	)
	mprincipal.Refresh()
}

func pushPR(pr gitter.PullRequest) {

	status, ok := decisions[pr.ReviewDecision]
	if !ok {
		status = ""
	}
	val := fmt.Sprintf("(#%d) %s\n[%s] ↦ [%s]\n%s", pr.Number, pr.Title, pr.HeadRefName, pr.BaseRefName, status)
	pushPRItem(val, func() {
		actions.OpenLink(pr.Permalink, Contexts[currentContext].ChromeProfile)
	})
}

func addItems() {
	prBox = fyne.NewMenuItem("PRs", func() {})
	prBox.ChildMenu = fyne.NewMenu("PRs")

	mprincipal.Items = []*fyne.MenuItem{
		fyne.NewMenuItem("change icon", func() {
			fmt.Println("Clicked")
			pushPRItem("Test PR", func() {})
			// desk.SetSystemTrayIcon(resIconReviewable)
		}),
		fyne.NewMenuItem("delete self", func() {
			fmt.Println("Clicked")
			mprincipal.Items = mprincipal.Items[:1] // how to delete stuff
			mprincipal.Refresh()
		}),
		fyne.NewMenuItemSeparator(),
		prBox,
	}
}

func ehp(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
}
