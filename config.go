package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadConfig() {
	file_contents, err := os.ReadFile(CONFIG_FILE)
	fmt.Println(file_contents)
	err = json.Unmarshal(file_contents, &Config)
	ehp(err)
}
