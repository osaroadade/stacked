package stack

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type StackEntry struct {
	Parent string `yaml:"parent"`
	PR     string `yaml:"pr,omitempty"`
	Status string `yaml:"status"`
}

type StackData struct {
	Stack map[string]StackEntry `yaml:"stack"`
}

func WriteSampleStack() error {
	data := StackData{
		Stack: map[string]StackEntry{
			"feature/stack-tracking": {
				Parent: "feature/base",
				Status: "open",
			},
		},
	}

	out, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(".stack.yaml", out, 0644)
	if err != nil {
		return err
	}

	fmt.Println("📝 .stack.yml written successfully.")
	return nil
}
