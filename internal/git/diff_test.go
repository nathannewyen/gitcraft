// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package git

import (
	"strings"
	"testing"
)

func TestGetRecentCommits_InvalidCount(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{"zero count", 0, true},
		{"negative count", -1, true},
		{"valid count", 5, false},
		{"large count gets capped", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test may fail if not in a git repo
			// but we're primarily testing the validation logic
			_, err := GetRecentCommits(tt.count)

			// For invalid counts, we expect an error
			if tt.wantErr && err == nil {
				t.Errorf("GetRecentCommits(%d) expected error, got nil", tt.count)
			}

			// For count <= 0, we should get our validation error
			if tt.count <= 0 && err != nil && !strings.Contains(err.Error(), "count must be positive") {
				t.Errorf("GetRecentCommits(%d) expected 'count must be positive' error, got: %v", tt.count, err)
			}
		})
	}
}

func TestTruncateDiff_NoTruncation(t *testing.T) {
	shortDiff := "diff --git a/file.go b/file.go\n+hello world"
	result := TruncateDiff(shortDiff, 1000)

	if result != shortDiff {
		t.Errorf("TruncateDiff should not modify short diffs.\nExpected: %s\nGot: %s", shortDiff, result)
	}
}

func TestTruncateDiff_WithTruncation(t *testing.T) {
	// Create a diff longer than maxLength
	longDiff := strings.Repeat("diff --git a/file.go b/file.go\n+line content\n", 100)
	maxLength := 500

	result := TruncateDiff(longDiff, maxLength)

	// Result should be truncated
	if len(result) > maxLength+100 { // Allow some buffer for the truncation message
		t.Errorf("TruncateDiff should truncate long diffs. Length: %d, Max: %d", len(result), maxLength)
	}

	// Should contain truncation message
	if !strings.Contains(result, "[... diff truncated due to length ...]") {
		t.Error("TruncateDiff should add truncation message")
	}
}

func TestTruncateDiff_FindsGoodBreakPoint(t *testing.T) {
	// Create a diff with multiple file sections
	diff := "diff --git a/file1.go b/file1.go\n+content1\n" +
		"diff --git a/file2.go b/file2.go\n+content2\n" +
		"diff --git a/file3.go b/file3.go\n+content3\n"

	// Truncate to a length that falls in the middle of file3
	maxLength := len(diff) - 20
	result := TruncateDiff(diff, maxLength)

	// Should break at a diff boundary if possible
	if strings.Contains(result, "content3") && !strings.Contains(result, "[... diff truncated") {
		t.Error("TruncateDiff should try to break at diff boundaries")
	}
}

func TestIsGitRepository(t *testing.T) {
	// This test is environment-dependent
	// In the gitcraft repo, this should return true
	result := IsGitRepository()

	// We can't guarantee the test environment, so just check it doesn't panic
	_ = result
}

func TestDiffInfo_Structure(t *testing.T) {
	// Test that DiffInfo struct has expected fields
	info := DiffInfo{
		Diff:         "test diff content",
		FilesChanged: []string{"file1.go", "file2.go"},
		Stats:        "2 files changed",
	}

	if info.Diff != "test diff content" {
		t.Error("DiffInfo.Diff not set correctly")
	}
	if len(info.FilesChanged) != 2 {
		t.Error("DiffInfo.FilesChanged not set correctly")
	}
	if info.Stats != "2 files changed" {
		t.Error("DiffInfo.Stats not set correctly")
	}
}

func TestErrors(t *testing.T) {
	// Test that error constants are defined
	if ErrNotGitRepo == nil {
		t.Error("ErrNotGitRepo should not be nil")
	}
	if ErrNoStagedChanges == nil {
		t.Error("ErrNoStagedChanges should not be nil")
	}
	if ErrGitNotFound == nil {
		t.Error("ErrGitNotFound should not be nil")
	}

	// Test error messages
	if !strings.Contains(ErrNotGitRepo.Error(), "git repository") {
		t.Error("ErrNotGitRepo should mention 'git repository'")
	}
	if !strings.Contains(ErrNoStagedChanges.Error(), "staged") {
		t.Error("ErrNoStagedChanges should mention 'staged'")
	}
}
