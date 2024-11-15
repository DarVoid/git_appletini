package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadTrackingConfig() {
	file_contents, err := os.ReadFile(TRACKING_CONFIG_FILE)
	if err != nil {
		fmt.Println()
	}
	// fmt.Println(string(file_contents))
	err = json.Unmarshal(file_contents, &TrackingConfig)
	ehp(err)
	fmt.Println(TrackingConfig)
}
