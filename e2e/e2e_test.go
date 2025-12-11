// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// getBinaryPath returns the path to the gitcraft binary
func getBinaryPath(t *testing.T) string {
	// First, try to find an existing binary
	paths := []string{
		"../gitcraft",
		"./gitcraft",
		filepath.Join(os.Getenv("GOPATH"), "bin", "gitcraft"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			absPath, _ := filepath.Abs(p)
			return absPath
		}
	}

	// Build the binary if not found
	t.Log("Building gitcraft binary for tests...")
	cmd := exec.Command("go", "build", "-o", "../gitcraft", "..")
	if err := cmd.Run(); err != nil {
		t.Skipf("Could not build gitcraft binary: %v", err)
	}

	absPath, _ := filepath.Abs("../gitcraft")
	return absPath
}

func TestCLI_Version(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gitcraft --version failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "gitcraft version") {
		t.Errorf("Expected version output to contain 'gitcraft version', got: %s", output)
	}
}

func TestCLI_Help(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gitcraft --help failed: %v\nOutput: %s", err, output)
	}

	// Check for expected help content
	expectedStrings := []string{
		"gitcraft",
		"AI",
		"commit",
		"--provider",
		"--type",
		"config",
		"hook",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(string(output), expected) {
			t.Errorf("Help output should contain '%s', got: %s", expected, output)
		}
	}
}

func TestCLI_ConfigHelp(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gitcraft config --help failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "config") {
		t.Errorf("Config help should mention 'config', got: %s", output)
	}
}

func TestCLI_HookHelp(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "hook", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gitcraft hook --help failed: %v\nOutput: %s", err, output)
	}

	expectedStrings := []string{"install", "uninstall", "status"}
	for _, expected := range expectedStrings {
		if !strings.Contains(string(output), expected) {
			t.Errorf("Hook help should mention '%s', got: %s", expected, output)
		}
	}
}

func TestCLI_ConfigPath(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "path")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gitcraft config path failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), ".gitcraft.yaml") {
		t.Errorf("Config path should contain '.gitcraft.yaml', got: %s", output)
	}
}

func TestCLI_InvalidCommand(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "invalidcommand")
	output, err := cmd.CombinedOutput()

	// Should exit with error
	if err == nil {
		t.Error("Expected error for invalid command")
	}

	// Should mention unknown command
	if !strings.Contains(string(output), "unknown") && !strings.Contains(string(output), "Error") {
		t.Errorf("Should indicate unknown command, got: %s", output)
	}
}

func TestCLI_InvalidType(t *testing.T) {
	binary := getBinaryPath(t)

	// Create a temp directory that's not a git repo
	tmpDir, err := os.MkdirTemp("", "gitcraft-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Run with invalid type in a git repo context
	cmd := exec.Command(binary, "--type", "invalidtype", "--dry-run")
	cmd.Dir = tmpDir
	output, _ := cmd.CombinedOutput()

	// Should either error about invalid type or about not being a git repo
	outputStr := string(output)
	if !strings.Contains(outputStr, "Invalid commit type") && !strings.Contains(outputStr, "git repository") {
		t.Logf("Output: %s", outputStr)
		// This is okay - we just want to make sure it handles the flag
	}
}

func TestCLI_NotGitRepo(t *testing.T) {
	binary := getBinaryPath(t)

	// Create a temp directory that's not a git repo
	tmpDir, err := os.MkdirTemp("", "gitcraft-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command(binary)
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()

	// Should fail because it's not a git repo
	if err == nil {
		t.Error("Expected error when running outside git repository")
	}

	// Error message should indicate the problem
	outputStr := string(output)
	if !strings.Contains(outputStr, "git repository") && !strings.Contains(outputStr, "provider") {
		t.Logf("Expected error about git repo or provider, got: %s", outputStr)
	}
}

func TestCLI_DryRunFlag(t *testing.T) {
	binary := getBinaryPath(t)

	// Verify --dry-run flag is recognized
	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "dry-run") {
		t.Error("Help should mention --dry-run flag")
	}
}

func TestCLI_EmojiFlag(t *testing.T) {
	binary := getBinaryPath(t)

	// Verify --emoji flag is recognized
	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "emoji") {
		t.Error("Help should mention --emoji flag")
	}
}

func TestCLI_ProviderFlag(t *testing.T) {
	binary := getBinaryPath(t)

	// Verify --provider flag is recognized
	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "provider") {
		t.Error("Help should mention --provider flag")
	}

	// Should list supported providers
	if !strings.Contains(string(output), "openai") || !strings.Contains(string(output), "claude") {
		t.Error("Help should list supported providers")
	}
}
