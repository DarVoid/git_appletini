package main

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"git_applet/actions"
	"git_applet/gitter"
	_ "git_applet/gitter"
	"git_applet/types"
	"net/http"
	"sync"
	"time"

	"os"

	"fyne.io/systray"
	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2"
)

//go:embed config.json
var b []byte

var data types.Config
var Contexts types.ContextMap
var currentContext string
var prBox *systray.MenuItem
var currentHash string = ""
var client *http.Client
var prs []gitter.PullRequest
var status *systray.MenuItem

//go:embed assets/tray1.png
var iconDefault []byte

//go:embed assets/tray2.png
var iconReviewable []byte

//go:embed assets/tray3.png
var iconMergeable []byte

//go:embed assets/tray4.png
var iconBoth []byte

var wg sync.WaitGroup

func openGithub(item *systray.MenuItem) {
	for range item.ClickedCh {
		url := fmt.Sprintf("https://%v", Contexts[currentContext].Github.Host)
		profile := fmt.Sprintf("%v", Contexts[currentContext].ChromeProfile)

		actions.OpenLink(url, profile)
	}
}
func handleExit(item *systray.MenuItem) {
	for range item.ClickedCh {
		systray.Quit()
	}
}

func handleIconChange(item *systray.MenuItem, data []byte) {
	for range item.ClickedCh {
		setTrayIcon(data)
	}
}

func setTrayIcon(data []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	systray.SetTemplateIcon(data, data)
}

func chooseContext(context string) {
	currentContext = context
	name := Contexts[currentContext].Github.Name
	username := Contexts[currentContext].Github.Username
	email := Contexts[currentContext].Github.Email
	host := Contexts[currentContext].Github.Host
	actions.ChangeToProfile(name, email, username, host)
	status.ClickedCh <- struct{}{}
}
func changeToContext(item *systray.MenuItem, context string) { // ,  ...item *systray.MenuItem
	for range item.ClickedCh {
		// parent := item.GetParent()
		// for _, sibling := range parent.GetChildren() {
		// 	sibling.Uncheck()
		// }
		// item.Check()
		// chooseContext(context)
	}
}

func addPolling(parent *systray.MenuItem, pollConf types.PollConfig, identifier string) {
	fmt.Printf("pollConf: %v\n", pollConf)
	go func() {
		for {
			if pollConf.Enabled {

				//TODO: Send a msg on the parent element to refresh the childs
				time.Sleep(time.Duration(pollConf.Frequency) * time.Second)
				fmt.Printf("%v: %v\n", Contexts[identifier].Github.Host, Contexts[identifier].Poll.Frequency)
				polledPRs()
			}
		}
	}()
}

func addContextChange(parent *systray.MenuItem, contextKey string, context types.Context) {
	novo := parent.AddSubMenuItemCheckbox(contextKey, contextKey, (data.DefaultContext == contextKey))
	go changeToContext(novo, contextKey)

	go addPolling(parent, context.Poll, contextKey)
}

func addContexOptions() {
	item := systray.AddMenuItem("Open Github", "Github")
	go openGithub(item)

	itemProfiles := systray.AddMenuItem("Context", "personal")

	for key, value := range data.Contexts {

		//##################################################################//
		//                          (Un)comentable                          //
		//##################################################################//
		// val, err := json.Marshal(types.Context{
		// 	Title:         value.Title,
		// 	ChromeProfile: value.ChromeProfile,
		// 	Github:        value.Github,
		// 	Poll:          value.Poll,
		// })
		// if err != nil {
		// 	fmt.Printf("err: %v\n", err)
		// }
		// fmt.Printf("value: %v\n", string(val))
		//##################################################################//
		//##################################################################//
		fmt.Printf("key: %v\n", key)
		addContextChange(itemProfiles, key, value)
		//##################################################################//
		//                          there's code here                       //
		//##################################################################//

	}

}

func exitOption() {
	itemLast := systray.AddMenuItem("Exit", "Bye")
	go handleExit(itemLast)
}

func separator() {
	systray.AddSeparator()
}

func handleDebug(item *systray.MenuItem) {
	for range item.ClickedCh {
		// fmt.Printf("item.GetChildren(): %v\n", item.GetChildren())
	}
}
func addPRContainer() {
	prBox = systray.AddMenuItem("PR Container", "container")
	return
}

func syncPolledItems() {
	decisions := make(map[string]string)
	decisions["APPROVED"] = "APPROVED ✔️"
	decisions["CHANGES_REQUESTED"] = "RIP, you tried💩"
	decisions[""] = "on Hold..."

	// prBox.RemoveChildren()
	green, red := false, false
	for _, pr := range prs {
		fmt.Printf("pr.ReviewRequests: %v\n", pr.ReviewRequests)
		if pr.ReviewDecision == "APPROVED" {
			green = true
		}
		if pr.ReviewDecision == "" {
			red = true
		}
		fmt.Printf("pr.ReviewDecision: %v\n", pr.ReviewDecision)
		status, ok := decisions[pr.ReviewDecision]

		if !ok {
			status = ""
		}
		if pr.ReviewCount > 0 {
			fmt.Printf("%v", "❗")
		}
		val := fmt.Sprintf("%v => %v\nStatus:%v\n\n -- %v", pr.HeadRefName, pr.BaseRefName, status, pr.Title)
		item := prBox.AddSubMenuItem(val, val)
		go handleLink(item, pr.Permalink)

	}
	if green && red {
		setTrayIcon(iconBoth)
		return
	}
	if red {
		setTrayIcon(iconReviewable)
		return
	}
	if green {
		setTrayIcon(iconMergeable)
		return
	}
	setTrayIcon(iconDefault)
}
func handleContextChangeStatusBar(item *systray.MenuItem) {
	for range item.ClickedCh {
		title := fmt.Sprintf("[%v] Connected 🌐", Contexts[currentContext].Title)
		item.SetTitle(title)
	}
}
func statusBar() (a *systray.MenuItem) {
	a = systray.AddMenuItem("", "")

	a.Disable()
	return
}
func startMain(channels ...chan struct{}) {
	for _, channel := range channels {
		channel <- struct{}{}
	}
}

func onReady() {

	setTrayIcon(iconDefault)
	wg.Add(1)
	status = statusBar()
	go handleContextChangeStatusBar(status)
	// iconChangeOptsExamples()
	addContexOptions()
	separator()

	// addSelfDeletingMenu()
	addPRContainer()
	go polledPRs()
	// addPolledItems(container)

	separator()

	exitOption()
	chooseContext(currentContext)
	fmt.Println("Running")
}

func onExit() {
	fmt.Println("Closing...")
}

func handleLink(item *systray.MenuItem, link string) {
	for range item.ClickedCh {
		actions.OpenLink(link, Contexts[currentContext].ChromeProfile)

	}
}

func polledPRs() {
	currentContext = data.DefaultContext
	Contexts = data.Contexts

	ctx := auth2()
	prsLocal := gitter.PrResponse{}
	gitter.GetPullRequests(Contexts[currentContext].Github.GraphQL, &prsLocal, os.Getenv(Contexts[currentContext].Github.Token), ctx)
	prs = prsLocal.Extract()
	cha := sha256.New()
	marco, err := json.Marshal(prs)
	ehp(err)
	cha.Write(marco)
	fmt.Printf("%x\n", string(cha.Sum(nil)))
	newHash := fmt.Sprintf("%x", string(cha.Sum(nil)))

	if currentHash != newHash {

		currentHash = newHash

		syncPolledItems()

	}

}

func main() {
	err := json.Unmarshal(b, &data)
	ehp(err)
	systray.Run(onReady, onExit)
}

func auth2() (ctx context.Context) {
	ctx = context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(Contexts[currentContext].Github.Token)},
	)
	client = oauth2.NewClient(ctx, ts)
	return
}
func ehp(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
}
