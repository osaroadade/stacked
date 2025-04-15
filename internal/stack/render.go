package stack

import (
	"fmt"
	"os"
	"sort"
)

// GenerateMarkdown creates a human-readable stack.md file based on .stack.yml.
func GenerateMarkdown(data StackData) error {
	var output string
	output += "# 🧱 Stack Overview\n\n"

	// Sort branches alphabetically for consistency
	branches := make([]string, 0, len(data.Stack))
	for branch := range data.Stack {
		branches = append(branches, branch)
	}
	sort.Strings(branches)

	for _, branch := range branches {
		entry := data.Stack[branch]

		// Status emoji
		statusEmoji := map[string]string{
			"merged": "✅",
			"open":   "🟡",
			"draft":  "📝",
		}[entry.Status]
		if statusEmoji == "" {
			statusEmoji = "❔"
		}

		// PR link (if available)
		prInfo := ""
		if entry.PR != nil {
			prInfo = fmt.Sprintf("[#%d](https://github.com/osaroadade/stacked/pull/%d)", *entry.PR, *entry.PR)
		}

		output += fmt.Sprintf("## %s %s\n", branch, statusEmoji)
		output += fmt.Sprintf("- Parent: %s\n", entry.Parent)
		if prInfo != "" {
			output += fmt.Sprintf("- PR: %s\n", prInfo)
		}
		output += "\n"
	}

	return os.WriteFile("stack.md", []byte(output), 0644)
}
