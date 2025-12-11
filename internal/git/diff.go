// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package git

import (
	"errors"
	"os/exec"
	"strings"
)

// Common errors for git operations
var (
	ErrNotGitRepo      = errors.New("not a git repository")
	ErrNoStagedChanges = errors.New("no staged changes found")
	ErrGitNotFound     = errors.New("git command not found")
)

// DiffInfo contains information about staged changes
type DiffInfo struct {
	// Diff is the raw diff output from git diff --staged
	Diff string
	// FilesChanged is a list of files that have been modified
	FilesChanged []string
	// Stats contains summary statistics (insertions, deletions)
	Stats string
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (*DiffInfo, error) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, ErrGitNotFound
	}

	// Check if we're in a git repository
	if !IsGitRepository() {
		return nil, ErrNotGitRepo
	}

	// Get the staged diff
	diffCmd := exec.Command("git", "diff", "--staged")
	diffOutput, err := diffCmd.Output()
	if err != nil {
		return nil, err
	}

	diff := string(diffOutput)
	if strings.TrimSpace(diff) == "" {
		return nil, ErrNoStagedChanges
	}

	// Get list of staged files
	filesCmd := exec.Command("git", "diff", "--staged", "--name-only")
	filesOutput, err := filesCmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(filesOutput)), "\n")
	if len(files) == 1 && files[0] == "" {
		files = []string{}
	}

	// Get diff stats
	statsCmd := exec.Command("git", "diff", "--staged", "--stat")
	statsOutput, err := statsCmd.Output()
	if err != nil {
		return nil, err
	}

	return &DiffInfo{
		Diff:         diff,
		FilesChanged: files,
		Stats:        string(statsOutput),
	}, nil
}

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRecentCommits returns the most recent commit messages for context
func GetRecentCommits(count int) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", "-n", string(rune(count+'0')))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			commits = append(commits, line)
		}
	}
	return commits, nil
}

// TruncateDiff truncates a diff to a maximum length to avoid token limits
func TruncateDiff(diff string, maxLength int) string {
	if len(diff) <= maxLength {
		return diff
	}

	// Find a good breaking point (end of a file section)
	truncated := diff[:maxLength]
	lastDiffMarker := strings.LastIndex(truncated, "\ndiff --git")
	if lastDiffMarker > maxLength/2 {
		truncated = truncated[:lastDiffMarker]
	}

	return truncated + "\n\n[... diff truncated due to length ...]"
}
