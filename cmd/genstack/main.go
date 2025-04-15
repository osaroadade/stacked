package main

import (
	"fmt"
	"os"

	"github.com/osaroadade/stacked/internal/stack"
)

func main() {
	content, err := os.ReadFile(".stack.yaml")
	if err != nil {
		fmt.Println("❌ Could not read .stack.yaml:", err)
		os.Exit(1)
	}

	var data stack.StackData
	err = stack.UnmarshalStack(content, &data)
	if err != nil {
		fmt.Println("❌ Failed to parse .stack.yaml:", err)
		os.Exit(1)
	}

	err = stack.GenerateMarkdown(data)
	if err != nil {
		fmt.Println("❌ Failed to generate stack.md:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Generated .github/stack.md successfully.")
}
