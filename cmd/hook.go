// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Hook script content for prepare-commit-msg
const hookScript = `#!/bin/sh
# gitcraft prepare-commit-msg hook
# This hook generates commit messages using AI

# Only run if this is a new commit (not amend, merge, etc.)
COMMIT_SOURCE=$2
if [ "$COMMIT_SOURCE" = "message" ] || [ "$COMMIT_SOURCE" = "merge" ] || [ "$COMMIT_SOURCE" = "squash" ]; then
    exit 0
fi

# Check if there are staged changes
if git diff --cached --quiet; then
    exit 0
fi

# Check if gitcraft is installed
if ! command -v gitcraft &> /dev/null; then
    echo "gitcraft not found, skipping AI commit message generation"
    exit 0
fi

# Generate the commit message
echo "Generating commit message with gitcraft..."
COMMIT_MSG=$(gitcraft --dry-run 2>/dev/null | grep -A 100 "Generated commit message:" | tail -n +2)

if [ -n "$COMMIT_MSG" ]; then
    echo "$COMMIT_MSG" > "$1"
fi
`

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git hooks for gitcraft",
	Long: `Manage git hooks to automatically generate commit messages.

The prepare-commit-msg hook will run gitcraft automatically
when you use 'git commit' without a message.

Examples:
  gitcraft hook install    Install the hook in current repo
  gitcraft hook uninstall  Remove the hook from current repo
  gitcraft hook status     Check if hook is installed`,
}

// hookInstallCmd installs the git hook
var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git hook in current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find git directory
		gitDir, err := getGitDir()
		if err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render("Not a git repository"))
		}

		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")

		// Check if hook already exists
		if _, err := os.Stat(hookPath); err == nil {
			// Read existing hook to check if it's ours
			content, err := os.ReadFile(hookPath)
			if err != nil {
				return fmt.Errorf("%s", ErrorStyle.Render(fmt.Sprintf("Failed to read existing hook: %v", err)))
			}

			if string(content) == hookScript {
				fmt.Println(InfoStyle.Render("gitcraft hook is already installed"))
				return nil
			}

			// Backup existing hook
			backupPath := hookPath + ".backup"
			if err := os.Rename(hookPath, backupPath); err != nil {
				return fmt.Errorf("%s", ErrorStyle.Render(fmt.Sprintf("Failed to backup existing hook: %v", err)))
			}
			fmt.Println(DimStyle.Render(fmt.Sprintf("Backed up existing hook to %s", backupPath)))
		}

		// Create hooks directory if it doesn't exist
		hooksDir := filepath.Join(gitDir, "hooks")
		if err := os.MkdirAll(hooksDir, 0755); err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render(fmt.Sprintf("Failed to create hooks directory: %v", err)))
		}

		// Write hook script
		if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render(fmt.Sprintf("Failed to write hook: %v", err)))
		}

		fmt.Println(SuccessStyle.Render("gitcraft hook installed successfully!"))
		fmt.Println(DimStyle.Render("The hook will generate commit messages when you run 'git commit'"))
		return nil
	},
}

// hookUninstallCmd removes the git hook
var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove git hook from current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitDir, err := getGitDir()
		if err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render("Not a git repository"))
		}

		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")

		// Check if hook exists
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			fmt.Println(DimStyle.Render("No hook to remove"))
			return nil
		}

		// Read hook to verify it's ours
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf(ErrorStyle.Render("Failed to read hook: %v"), err)
		}

		if string(content) != hookScript {
			fmt.Println(ErrorStyle.Render("Warning: The current hook was not installed by gitcraft"))
			fmt.Print("Remove anyway? [y/N]: ")

			var input string
			fmt.Scanln(&input)
			if input != "y" && input != "Y" {
				fmt.Println(DimStyle.Render("Cancelled"))
				return nil
			}
		}

		// Remove hook
		if err := os.Remove(hookPath); err != nil {
			return fmt.Errorf(ErrorStyle.Render("Failed to remove hook: %v"), err)
		}

		// Restore backup if exists
		backupPath := hookPath + ".backup"
		if _, err := os.Stat(backupPath); err == nil {
			if err := os.Rename(backupPath, hookPath); err != nil {
				fmt.Println(DimStyle.Render(fmt.Sprintf("Note: Failed to restore backup hook: %v", err)))
			} else {
				fmt.Println(DimStyle.Render("Restored previous hook from backup"))
			}
		}

		fmt.Println(SuccessStyle.Render("gitcraft hook removed successfully!"))
		return nil
	},
}

// hookStatusCmd checks the hook status
var hookStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check git hook status",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitDir, err := getGitDir()
		if err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render("Not a git repository"))
		}

		hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")

		// Check if hook exists
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			fmt.Println(DimStyle.Render("gitcraft hook is not installed"))
			fmt.Println(DimStyle.Render("Run 'gitcraft hook install' to install it"))
			return nil
		}

		// Read hook to verify it's ours
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("%s", ErrorStyle.Render(fmt.Sprintf("Failed to read hook: %v", err)))
		}

		if string(content) == hookScript {
			fmt.Println(SuccessStyle.Render("gitcraft hook is installed"))
		} else {
			fmt.Println(InfoStyle.Render("A different prepare-commit-msg hook is installed"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
	hookCmd.AddCommand(hookStatusCmd)
}

// getGitDir returns the path to the .git directory
func getGitDir() (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := gitCmd.Output()
	if err != nil {
		return "", err
	}
	// Safety: Check length before slicing to prevent index out of bounds panic
	if len(output) == 0 {
		return "", fmt.Errorf("empty output from git rev-parse")
	}
	// Remove trailing newline if present
	result := string(output)
	if result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return filepath.Clean(result), nil
}
