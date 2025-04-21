package stack

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type StackEntry struct {
	Parent string `yaml:"parent"`
	PR     *int   `yaml:"pr,omitempty"`
	Status string `yaml:"status"`
}

type StackData struct {
	Stack map[string]StackEntry `yaml:"stack"`
}

func WriteBranchEntry(branch, parent string, pr int) error {
	entry := StackEntry{
		Parent: parent,
		PR:     &pr,
		Status: "open",
	}

	// Load existing file if it exists
	var data StackData
	content, err := os.ReadFile(".stack.yaml")
	if err == nil {
		yaml.Unmarshal(content, &data)
	} else {
		data = StackData{
			Stack: make(map[string]StackEntry),
		}
	}

	// Update/add entry
	data.Stack[branch] = entry

	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(".stack.yaml", out, 0644)
	if err != nil {
		return err
	}

	fmt.Println("📝 .stack.yml updated successfully.")
	return GenerateMarkdown(data)
}
