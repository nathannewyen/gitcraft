// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package ai

import (
	"strings"
	"testing"
)

func TestBuildPrompt_BasicDiff(t *testing.T) {
	req := &GenerateRequest{
		Diff:         "+hello world",
		FilesChanged: []string{"main.go"},
		MaxLength:    72,
	}

	prompt := BuildPrompt(req)

	// Should contain essential instructions
	if !strings.Contains(prompt, "Conventional Commits") {
		t.Error("Prompt should mention Conventional Commits")
	}
	if !strings.Contains(prompt, "72") {
		t.Error("Prompt should include max length")
	}
	if !strings.Contains(prompt, "main.go") {
		t.Error("Prompt should include changed files")
	}
	if !strings.Contains(prompt, "+hello world") {
		t.Error("Prompt should include the diff")
	}
}

func TestBuildPrompt_WithCommitType(t *testing.T) {
	req := &GenerateRequest{
		Diff:       "+new feature",
		CommitType: "feat",
		MaxLength:  72,
	}

	prompt := BuildPrompt(req)

	if !strings.Contains(prompt, "Commit type: feat") {
		t.Error("Prompt should include specified commit type")
	}
}

func TestBuildPrompt_WithEmoji(t *testing.T) {
	req := &GenerateRequest{
		Diff:      "+emoji test",
		UseEmoji:  true,
		MaxLength: 72,
	}

	prompt := BuildPrompt(req)

	if !strings.Contains(prompt, "gitmoji") {
		t.Error("Prompt should mention gitmoji when UseEmoji is true")
	}
}

func TestBuildPrompt_WithRecentCommits(t *testing.T) {
	req := &GenerateRequest{
		Diff:          "+test",
		RecentCommits: []string{"abc123 feat: previous feature", "def456 fix: bug fix"},
		MaxLength:     72,
	}

	prompt := BuildPrompt(req)

	if !strings.Contains(prompt, "Recent commits") {
		t.Error("Prompt should include recent commits section")
	}
	if !strings.Contains(prompt, "previous feature") {
		t.Error("Prompt should include recent commit content")
	}
}

func TestParseResponse_SubjectOnly(t *testing.T) {
	response := "feat: add new feature"

	result := ParseResponse(response)

	if result.Subject != "feat: add new feature" {
		t.Errorf("Expected subject 'feat: add new feature', got '%s'", result.Subject)
	}
	if result.Body != "" {
		t.Errorf("Expected empty body, got '%s'", result.Body)
	}
}

func TestParseResponse_WithBody(t *testing.T) {
	response := `feat: add new feature

This commit adds a new feature that does something useful.
It includes multiple lines of description.`

	result := ParseResponse(response)

	if result.Subject != "feat: add new feature" {
		t.Errorf("Expected subject 'feat: add new feature', got '%s'", result.Subject)
	}
	if !strings.Contains(result.Body, "new feature that does something useful") {
		t.Errorf("Body should contain description, got '%s'", result.Body)
	}
}

func TestParseResponse_WithCodeBlocks(t *testing.T) {
	response := "```\nfeat: add feature\n```"

	result := ParseResponse(response)

	// Should strip code blocks
	if strings.Contains(result.Subject, "```") {
		t.Error("ParseResponse should strip markdown code blocks")
	}
	if result.Subject != "feat: add feature" {
		t.Errorf("Expected 'feat: add feature', got '%s'", result.Subject)
	}
}

func TestParseResponse_WithWhitespace(t *testing.T) {
	response := "   fix: bug fix   \n\n"

	result := ParseResponse(response)

	if result.Subject != "fix: bug fix" {
		t.Errorf("Expected trimmed subject 'fix: bug fix', got '%s'", result.Subject)
	}
}

func TestParseResponse_EmptyResponse(t *testing.T) {
	response := ""

	result := ParseResponse(response)

	if result.Subject != "" {
		t.Errorf("Expected empty subject for empty response, got '%s'", result.Subject)
	}
}

func TestGenerateRequest_Structure(t *testing.T) {
	req := &GenerateRequest{
		Diff:          "test diff",
		FilesChanged:  []string{"file1.go", "file2.go"},
		CommitType:    "feat",
		UseEmoji:      true,
		MaxLength:     50,
		RecentCommits: []string{"commit1", "commit2"},
	}

	if req.Diff != "test diff" {
		t.Error("Diff not set correctly")
	}
	if len(req.FilesChanged) != 2 {
		t.Error("FilesChanged not set correctly")
	}
	if req.CommitType != "feat" {
		t.Error("CommitType not set correctly")
	}
	if !req.UseEmoji {
		t.Error("UseEmoji not set correctly")
	}
	if req.MaxLength != 50 {
		t.Error("MaxLength not set correctly")
	}
}

func TestGenerateResponse_Structure(t *testing.T) {
	resp := &GenerateResponse{
		Subject: "feat: test",
		Body:    "Test body",
	}

	if resp.Subject != "feat: test" {
		t.Error("Subject not set correctly")
	}
	if resp.Body != "Test body" {
		t.Error("Body not set correctly")
	}
}

func TestErrors_Defined(t *testing.T) {
	if ErrNoAPIKey == nil {
		t.Error("ErrNoAPIKey should not be nil")
	}
	if ErrProviderNotFound == nil {
		t.Error("ErrProviderNotFound should not be nil")
	}
	if ErrGenerationFailed == nil {
		t.Error("ErrGenerationFailed should not be nil")
	}
}
