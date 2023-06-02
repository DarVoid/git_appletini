package actions

import (
	"fmt"
	"os/exec"
)

func OpenLink(url string, profile string) {
	a := exec.Command("google-chrome", url, fmt.Sprintf("--profile-directory=%v", profile))
	a.Start()
}

func ChangeToProfile(name, email, username, host string) {

	a := exec.Command("git", "config", "--global", "user.name", fmt.Sprintf("\"%v\"", name)) //"${userName}"`
	a.Run()
	a = exec.Command("git", "config", "--global", "user.email", fmt.Sprintf("\"%v\"", email))
	a.Run()
	a = exec.Command("git", "config", "--global", "credential.helper", "store")
	a.Run()
	a = exec.Command("git", "config", "--global", fmt.Sprintf("credential.https://%v.username", host), username)
	a.Run()

}
