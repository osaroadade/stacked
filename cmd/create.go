/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/osaroadade/stacked/internal/stack"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a PR fro the current stacked branch",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Detect current branch
		output, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			fmt.Println("❌ Failed to get current branch:", err)
			return
		}

		currentBranch := strings.TrimSpace(string(output))
		fmt.Println("🌿 Current branch:", currentBranch)

		// Try to detect parent branch
		parentBranch, err := findParentBranch(currentBranch)
		if err != nil {
			fmt.Println("⚠️ Could not determine parent branch:", err)
		} else {
			fmt.Println("🔗 Detected parent branch:", parentBranch)
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("📝 Enter PR title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Print("📝 Enter PR description (optional): ")
		body, _ := reader.ReadString('\n')
		body = strings.TrimSpace(body)

		repoURL, err := getGitHubRepoURL()
		if err != nil {
			fmt.Println("⚠️ Could not get remote repo URL:", err)
			repoURL = "https://github.com/unkown/repo" // fallback
		}

		stackLink := fmt.Sprintf(
			"\n\n---\n🔗 This PR is part of a stack. See full context: [stack.md](%s/blob/%s/stack.md)",
			repoURL, currentBranch,
		)

		var fullBody string
		if body == "" {
			fullBody = stackLink
		} else {
			fullBody = body + stackLink
		}

		cmdArgs := []string{
			"pr", "create",
			"--title", title,
			"--body", fullBody,
			"--base", parentBranch,
			"--head", currentBranch,
		}

		fmt.Println("📤 Creating PR via GitHub CLI...")
		createCmd := exec.Command("gh", cmdArgs...)
		var prOutput bytes.Buffer
		createCmd.Stdout = &prOutput
		createCmd.Stderr = os.Stderr

		if err := createCmd.Run(); err != nil {
			fmt.Println("❌ Failed to create PR:", err)
			return
		}

		if err := createCmd.Run(); err != nil {
			fmt.Println("❌ Failed to create PR:", err)
			return
		}

		prURL := strings.TrimSpace(prOutput.String())
		fmt.Println("✅ PR created successfully!", prURL)

		re := regexp.MustCompile(`/pull/(\d+)$`)
		matches := re.FindStringSubmatch(prURL)
		if len(matches) < 2 {
			fmt.Println("⚠️ Could not extract PR number from URL:", prURL)
			return
		}
		prNumber, _ := strconv.Atoi(matches[1])

		err = stack.WriteSampleStack(currentBranch, parentBranch, prNumber)
		if err != nil {
			fmt.Println("⚠️ Could not write to .stack.yml:", err)
		}
	},
}

func init() {
	prCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func findParentBranch(currentBranch string) (string, error) {
	// Get all local branches
	branchesOut, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return "", err
	}
	allBranches := strings.Split(strings.TrimSpace(string(branchesOut)), "\n")

	// Try to find the best parent
	var bestParent string
	var bestDistance int

	for _, branch := range allBranches {
		branch = strings.TrimSpace(branch)
		if branch == currentBranch {
			continue
		}

		// Run: git merge-base currentBranch otherBranch
		mbCmd := exec.Command("git", "merge-base", currentBranch, branch)
		mergeBase, err := mbCmd.Output()
		if err != nil {
			continue // Skip if no common base
		}

		// Count commits between merge base and current branch
		countCmd := exec.Command("git", "rev-list", "--count", currentBranch+"^@", "^"+strings.TrimSpace(string(mergeBase)))
		countOut, err := countCmd.Output()
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(strings.TrimSpace(string(countOut)))
		if err != nil {
			continue
		}

		if bestParent == "" || count < bestDistance {
			bestParent = branch
			bestDistance = count
		}
	}

	if bestParent == "" {
		return "", fmt.Errorf("could not determine parent branch")
	}

	return bestParent, nil
}

func getGitHubRepoURL() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}
	rawURL := strings.TrimSpace(string(out))

	// Convert SSH URL to HTTPS
	if strings.HasPrefix(rawURL, "git@") {
		// Example: git@github.com:user/repo.git
		rawURL = strings.Replace(rawURL, "git@", "https://", 1)
		rawURL = strings.Replace(rawURL, ":", "/", 1)
	} else if strings.HasPrefix(rawURL, "https://") {
		// already done
	} else {
		return "", fmt.Errorf("unsupported remote URL format")
	}

	// Trim .git suffix
	if strings.HasSuffix(rawURL, ".git") {
		rawURL = strings.TrimSuffix(rawURL, ".git")
	}

	return rawURL, nil
}
