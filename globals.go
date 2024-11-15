package main

import (
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"

	"git_applet/gitter"
	"git_applet/types"
)

// global and default stuff
const CONFIG_FILE = "config.json"

var currentContext string
var Contexts types.ContextMap
var Config types.Config
var currentHash string = ""
var prs []gitter.PullRequest

// tracked configured stuff
const TRACKING_CONFIG_FILE = "tracking.config.json"

var TrackingConfig types.TrackingConfig
var PrQuery io.Writer
var SavedPRQuerry string
var currentTrackedHash string = ""
var trackedPrs []gitter.PullRequest

// visual stuff
var prBox *fyne.MenuItem
var trackedPrBox *fyne.MenuItem
var mprincipal *fyne.Menu
var contextSelector *fyne.MenuItem
var status *systray.MenuItem
var desk desktop.App

// global singleton
var client *http.Client
