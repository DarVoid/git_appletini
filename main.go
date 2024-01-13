package main

import (
	_ "embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	a := app.New()

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Git Applet")

		desk.SetSystemTrayMenu(m)

		a.Lifecycle().SetOnStarted(func() {
			desk.SetSystemTrayIcon(resIconDefault)
		})

		m.Items = []*fyne.MenuItem{
			fyne.NewMenuItem("Show", func() {
				fmt.Println("Clicked")
				desk.SetSystemTrayIcon(resIconReviewable)
			}),
			fyne.NewMenuItem("Show 2", func() {
				fmt.Println("Clicked")
				m.Items = m.Items[:1] // how to delete stuff
				m.Refresh()
			}),
		}
	}
	a.Run()
}
