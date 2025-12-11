// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Common errors for AI operations
var (
	ErrNoAPIKey         = errors.New("API key not configured")
	ErrProviderNotFound = errors.New("provider not found")
	ErrGenerationFailed = errors.New("failed to generate commit message")
)

// GenerateRequest contains the parameters for generating a commit message
type GenerateRequest struct {
	// Diff is the git diff content
	Diff string
	// FilesChanged is the list of changed files
	FilesChanged []string
	// CommitType is the optional commit type (feat, fix, etc.)
	CommitType string
	// UseEmoji indicates whether to use gitmoji
	UseEmoji bool
	// MaxLength is the maximum length for the subject line
	MaxLength int
	// RecentCommits are recent commit messages for style reference
	RecentCommits []string
}

// GenerateResponse contains the generated commit message
type GenerateResponse struct {
	// Subject is the commit message subject line
	Subject string
	// Body is the optional commit message body
	Body string
}

// Provider is the interface that all AI providers must implement
type Provider interface {
	// Name returns the name of the provider
	Name() string
	// Generate generates a commit message based on the request
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	// IsConfigured returns true if the provider is properly configured
	IsConfigured() bool
}

// BuildPrompt creates the prompt for AI providers
func BuildPrompt(req *GenerateRequest) string {
	var sb strings.Builder

	// System instruction
	sb.WriteString("You are an expert at writing concise, meaningful git commit messages.\n")
	sb.WriteString("Generate a commit message following the Conventional Commits specification.\n\n")

	// Format rules
	sb.WriteString("Rules:\n")
	sb.WriteString(fmt.Sprintf("- Subject line must be %d characters or less\n", req.MaxLength))
	sb.WriteString("- Use imperative mood (\"Add feature\" not \"Added feature\")\n")
	sb.WriteString("- Do not end the subject line with a period\n")
	sb.WriteString("- Separate subject from body with a blank line\n")
	sb.WriteString("- Body should explain what and why, not how\n\n")

	// Commit type
	if req.CommitType != "" {
		sb.WriteString(fmt.Sprintf("Commit type: %s\n\n", req.CommitType))
	} else {
		sb.WriteString("Determine the appropriate commit type from: feat, fix, docs, style, refactor, test, chore, perf, ci, build\n\n")
	}

	// Emoji
	if req.UseEmoji {
		sb.WriteString("Include a relevant gitmoji at the start of the subject line.\n")
		sb.WriteString("Common gitmojis: ✨ feat, 🐛 fix, 📝 docs, 💄 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf\n\n")
	}

	// Files changed
	if len(req.FilesChanged) > 0 {
		sb.WriteString("Files changed:\n")
		for _, file := range req.FilesChanged {
			sb.WriteString(fmt.Sprintf("- %s\n", file))
		}
		sb.WriteString("\n")
	}

	// Recent commits for style reference
	if len(req.RecentCommits) > 0 {
		sb.WriteString("Recent commits in this repository (for style reference):\n")
		for _, commit := range req.RecentCommits {
			sb.WriteString(fmt.Sprintf("- %s\n", commit))
		}
		sb.WriteString("\n")
	}

	// The diff
	sb.WriteString("Git diff:\n```\n")
	sb.WriteString(req.Diff)
	sb.WriteString("\n```\n\n")

	// Output format
	sb.WriteString("Generate a commit message. Output ONLY the commit message with no additional explanation.\n")
	sb.WriteString("Format:\n")
	sb.WriteString("<type>: <subject>\n\n")
	sb.WriteString("<optional body>\n")

	return sb.String()
}

// ParseResponse extracts subject and body from the AI response
func ParseResponse(response string) *GenerateResponse {
	response = strings.TrimSpace(response)

	// Remove any markdown code blocks
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Split into lines
	lines := strings.Split(response, "\n")
	if len(lines) == 0 {
		return &GenerateResponse{Subject: response}
	}

	subject := strings.TrimSpace(lines[0])
	var body string

	// Check for body (after blank line)
	if len(lines) > 2 {
		bodyLines := []string{}
		foundBlank := false
		for _, line := range lines[1:] {
			if strings.TrimSpace(line) == "" {
				foundBlank = true
				continue
			}
			if foundBlank {
				bodyLines = append(bodyLines, line)
			}
		}
		if len(bodyLines) > 0 {
			body = strings.Join(bodyLines, "\n")
		}
	}

	return &GenerateResponse{
		Subject: subject,
		Body:    body,
	}
}
