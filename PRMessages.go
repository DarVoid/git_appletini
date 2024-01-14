package main

const SHOW_RED_ON = "CHANGES_REQUESTED"
const SHOW_GREEN_ON = "APPROVED"

var decision_messages = map[string]string{
	"APPROVED":          "✅ Approved",
	"CHANGES_REQUESTED": "🚨 Changes requested",
	"REVIEW_REQUIRED":   "🛂 Requires review",
	"":                  "🕒 On hold...",
}

var merge_messages = map[string]string{
	"MERGEABLE":   "💚 Can be merged",
	"CONFLICTING": "🚩 Has conflicts",
}
