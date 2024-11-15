package main

import (
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"

	"git_applet/gitter"
	"git_applet/types"
)

const CONFIG_FILE = "config.json"

var Config types.Config
var Contexts types.ContextMap
var currentContext string
var prBox *fyne.MenuItem
var contextSelector *fyne.MenuItem
var currentHash string = ""
var client *http.Client
var prs []gitter.PullRequest
var status *systray.MenuItem
var mprincipal *fyne.Menu
var desk desktop.App
