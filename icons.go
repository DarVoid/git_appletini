package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
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
