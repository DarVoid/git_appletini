package main

import (
	"encoding/json"
	"os"
)

func loadConfig() {
	file_contents, err := os.ReadFile(CONFIG_FILE)
	err = json.Unmarshal(file_contents, &config)
	ehp(err)
}
