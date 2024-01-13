package main

import (
	_ "embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
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

	mprincipal.Items = []*fyne.MenuItem{
		fyne.NewMenuItem("Show", func() {
			fmt.Println("Clicked")
			desk.SetSystemTrayIcon(resIconReviewable)
		}),
		fyne.NewMenuItem("Show 2", func() {
			fmt.Println("Clicked")
			mprincipal.Items = mprincipal.Items[:1] // how to delete stuff
			mprincipal.Refresh()
		}),
	}
	go polledPRs()
	a.Run()
}

func ehp(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
}
