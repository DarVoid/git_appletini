package actions

import (
	"os/exec"
)

func OpenLink(url string, profile string) {
	a := exec.Command("open", url)
	a.Start()
}
