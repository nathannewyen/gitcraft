// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nathannewyen/gitcraft/internal/ai"
	"github.com/nathannewyen/gitcraft/internal/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Maximum diff size to send to AI (in characters)
const maxDiffSize = 8000

// generateCommitMessage is the main function that generates commit messages
func generateCommitMessage(cmd *cobra.Command) error {
	// Get flags
	commitType, _ := cmd.Flags().GetString("type")
	useEmoji, _ := cmd.Flags().GetBool("emoji")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	maxLength, _ := cmd.Flags().GetInt("max-length")

	// Get staged diff
	fmt.Println(InfoStyle.Render("Analyzing staged changes..."))

	diffInfo, err := git.GetStagedDiff()
	if err != nil {
		switch err {
		case git.ErrNotGitRepo:
			return fmt.Errorf(ErrorStyle.Render("Not a git repository"))
		case git.ErrNoStagedChanges:
			return fmt.Errorf(ErrorStyle.Render("No staged changes found. Use 'git add' to stage changes."))
		case git.ErrGitNotFound:
			return fmt.Errorf(ErrorStyle.Render("Git command not found. Please install git."))
		default:
			return fmt.Errorf(ErrorStyle.Render("Failed to get staged diff: %v"), err)
		}
	}

	// Show what we're working with
	fmt.Printf(DimStyle.Render("  Files changed: %d\n"), len(diffInfo.FilesChanged))
	for _, file := range diffInfo.FilesChanged {
		fmt.Printf(DimStyle.Render("    - %s\n"), file)
	}

	// Truncate diff if too large
	diff := git.TruncateDiff(diffInfo.Diff, maxDiffSize)

	// Get recent commits for style reference
	recentCommits, _ := git.GetRecentCommits(5)

	// Get provider
	provider, err := getProvider()
	if err != nil {
		return err
	}

	fmt.Printf(InfoStyle.Render("Generating commit message using %s...\n"), provider.Name())

	// Create request
	req := &ai.GenerateRequest{
		Diff:          diff,
		FilesChanged:  diffInfo.FilesChanged,
		CommitType:    commitType,
		UseEmoji:      useEmoji,
		MaxLength:     maxLength,
		RecentCommits: recentCommits,
	}

	// Generate with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := provider.Generate(ctx, req)
	if err != nil {
		return fmt.Errorf(ErrorStyle.Render("Failed to generate commit message: %v"), err)
	}

	// Display result
	fmt.Println()
	fmt.Println(SuccessStyle.Render("Generated commit message:"))
	fmt.Println()

	// Format the commit message nicely
	commitMsg := resp.Subject
	if resp.Body != "" {
		commitMsg += "\n\n" + resp.Body
	}

	// Display in a box
	displayCommitMessage(resp.Subject, resp.Body)

	if dryRun {
		fmt.Println()
		fmt.Println(DimStyle.Render("(dry-run mode - no commit created)"))
		return nil
	}

	// Ask user for confirmation
	fmt.Println()
	fmt.Print("Use this commit message? [Y/n/e(dit)]: ")

	var input string
	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "", "y", "yes":
		// Commit with the message
		return performCommit(commitMsg)
	case "e", "edit":
		// Open editor for editing
		return editAndCommit(commitMsg)
	default:
		fmt.Println(DimStyle.Render("Commit cancelled"))
		return nil
	}
}

// getProvider returns the appropriate AI provider based on configuration
func getProvider() (ai.Provider, error) {
	providerName := viper.GetString("provider")

	var provider ai.Provider

	switch providerName {
	case "openai":
		provider = ai.NewOpenAIProvider()
	case "claude":
		provider = ai.NewClaudeProvider()
	case "ollama":
		provider = ai.NewOllamaProvider()
	case "gemini":
		provider = ai.NewGeminiProvider()
	default:
		// Try to find a configured provider
		providers := []ai.Provider{
			ai.NewOpenAIProvider(),
			ai.NewClaudeProvider(),
			ai.NewOllamaProvider(),
			ai.NewGeminiProvider(),
		}

		for _, p := range providers {
			if p.IsConfigured() {
				provider = p
				break
			}
		}

		if provider == nil {
			return nil, fmt.Errorf(ErrorStyle.Render("No AI provider configured.\n") +
				DimStyle.Render("Run 'gitcraft config set <provider>.api_key <your-key>' to configure.\n") +
				DimStyle.Render("Supported providers: openai, claude, ollama, gemini"))
		}
	}

	if !provider.IsConfigured() {
		// Special case for Ollama which doesn't need API key
		if providerName != "ollama" {
			return nil, fmt.Errorf(ErrorStyle.Render("Provider '%s' not configured.\n")+
				DimStyle.Render("Run 'gitcraft config set %s.api_key <your-key>' to configure."),
				providerName, providerName)
		}
	}

	return provider, nil
}

// displayCommitMessage displays the commit message in a formatted box
func displayCommitMessage(subject, body string) {
	// Simple box display
	border := strings.Repeat("─", 60)

	fmt.Println("┌" + border + "┐")
	fmt.Printf("│ %s%s │\n", TitleStyle.Render(subject), strings.Repeat(" ", max(0, 58-len(subject))))

	if body != "" {
		fmt.Println("│" + strings.Repeat(" ", 60) + "│")
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			if len(line) > 58 {
				line = line[:55] + "..."
			}
			fmt.Printf("│ %s%s │\n", DimStyle.Render(line), strings.Repeat(" ", max(0, 58-len(line))))
		}
	}
	fmt.Println("└" + border + "┘")
}

// performCommit creates a git commit with the given message
func performCommit(message string) error {
	fmt.Println(SuccessStyle.Render("Committing..."))

	// Execute git commit using exec.Command
	gitCmd := exec.Command("git", "commit", "-m", message)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf(ErrorStyle.Render("Failed to commit: %v"), err)
	}

	fmt.Println(SuccessStyle.Render("Commit created successfully!"))
	return nil
}

// editAndCommit opens the user's editor to edit the commit message
func editAndCommit(message string) error {
	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	// Create a temporary file with the message
	tmpFile, err := os.CreateTemp("", "gitcraft-commit-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(message); err != nil {
		return fmt.Errorf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Open editor using exec.Command
	editorCmd := exec.Command(editor, tmpFile.Name())
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %v", err)
	}

	// Read edited message
	editedMsg, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read edited message: %v", err)
	}

	editedMsgStr := strings.TrimSpace(string(editedMsg))
	if editedMsgStr == "" {
		fmt.Println(DimStyle.Render("Empty message, commit cancelled"))
		return nil
	}

	return performCommit(editedMsgStr)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
