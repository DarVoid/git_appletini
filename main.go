package main

import (
	_ "embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

//go:embed assets/tray1.png
var iconDefault []byte
var resIconDefault = &fyne.StaticResource{
	StaticName:    "tray1.png",
	StaticContent: iconDefault,
}

//go:embed assets/tray2.png
var iconReviewable []byte
var resIconReviewable = &fyne.StaticResource{
	StaticName:    "tray2.png",
	StaticContent: iconReviewable,
}

//go:embed assets/tray3.png
var iconMergeable []byte
var resIconMergeable = &fyne.StaticResource{
	StaticName:    "tray3.png",
	StaticContent: iconMergeable,
}

//go:embed assets/tray4.png
var iconBoth []byte
var resIconBoth = &fyne.StaticResource{
	StaticName:    "tray4.png",
	StaticContent: iconBoth,
}

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
				m.Items = m.Items[:1]
				m.Refresh()
			}),
		}
	}
	a.Run()
}
