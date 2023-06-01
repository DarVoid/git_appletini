package actions

import (
	"fmt"
	"os/exec"
)

func OpenLink(url string, profile string) {
	a := exec.Command("google-chrome", url, fmt.Sprintf("--profile-directory:%v", profile))
	a.Start()
}
