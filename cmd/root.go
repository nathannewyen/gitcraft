// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version information set during build
var (
	Version   = "dev"
	BuildDate = "unknown"
)

// Style definitions for CLI output
var (
	// TitleStyle is used for the main title/banner
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED"))

	// SuccessStyle is used for successful operations
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	// InfoStyle is used for informational messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6"))

	// DimStyle is used for secondary/dimmed text
	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitcraft",
	Short: "AI-powered git commit message generator",
	Long: TitleStyle.Render(`
  ▄▀  █ ▀█▀ ▄▀▀ █▀▄ ▄▀▄ █▀ ▀█▀
  ▀▄█ █  █  ▀▄▄ █▀▄ █▀█ █▀  █
`) + `

Craft perfect git commit messages with AI.

gitcraft analyzes your staged changes and generates
meaningful commit messages following Conventional Commits.

` + DimStyle.Render("Supports: OpenAI, Claude, Ollama, Gemini"),
	Version: Version,
	RunE:    runGenerate,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringP("provider", "p", "", "AI provider to use (openai, claude, ollama, gemini)")
	rootCmd.PersistentFlags().StringP("model", "m", "", "Model to use for generation")

	// Flags for generate command (root command)
	rootCmd.Flags().StringP("type", "t", "", "Commit type (feat, fix, docs, style, refactor, test, chore)")
	rootCmd.Flags().BoolP("emoji", "e", false, "Add gitmoji to commit message")
	rootCmd.Flags().BoolP("dry-run", "d", false, "Show generated message without committing")
	rootCmd.Flags().IntP("max-length", "l", 72, "Maximum length for commit message subject")

	// Bind flags to viper
	viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider"))
	viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, ErrorStyle.Render("Error finding home directory:"), err)
		return
	}

	// Search config in home directory with name ".gitcraft"
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".gitcraft")

	// Read environment variables with GITCRAFT_ prefix
	viper.SetEnvPrefix("GITCRAFT")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("provider", "openai")
	viper.SetDefault("openai.model", "gpt-4o-mini")
	viper.SetDefault("claude.model", "claude-3-5-sonnet-20241022")
	viper.SetDefault("ollama.model", "llama3.2")
	viper.SetDefault("ollama.url", "http://localhost:11434")
	viper.SetDefault("gemini.model", "gemini-1.5-flash")

	// Try to read config file (ignore error if not found)
	if err := viper.ReadInConfig(); err == nil {
		// Config file found and successfully parsed
	}
}

// runGenerate is the main function that generates commit messages
func runGenerate(cmd *cobra.Command, args []string) error {
	// This will be implemented in generate.go
	return generateCommitMessage(cmd)
}
