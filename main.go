package main

import (
	_ "embed"
	"fmt"

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

func addItems() {
	prBox = fyne.NewMenuItem("PRs", func() {})

	// prBox.ChildMenu.Items = []*fyne.MenuItem{
	// 	fyne.NewMenuItem("Hello there", func() {}),
	// 	fyne.NewMenuItemSeparator(),
	// }

	mprincipal.Items = []*fyne.MenuItem{
		fyne.NewMenuItem("change icon", func() {
			fmt.Println("Clicked")
			desk.SetSystemTrayIcon(resIconReviewable)
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
