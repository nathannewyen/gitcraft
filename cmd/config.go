// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gitcraft configuration",
	Long: `Manage gitcraft configuration.

Configuration is stored in ~/.gitcraft.yaml

Examples:
  gitcraft config set openai.api_key sk-xxx
  gitcraft config set provider claude
  gitcraft config get provider
  gitcraft config list`,
}

// configSetCmd represents the config set subcommand
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		// Set the value in viper
		viper.Set(key, value)

		// Write to config file
		if err := writeConfig(); err != nil {
			return fmt.Errorf(ErrorStyle.Render("Failed to write config: %v"), err)
		}

		// Mask API keys in output
		displayValue := value
		if strings.Contains(strings.ToLower(key), "api_key") {
			if len(value) > 8 {
				displayValue = value[:4] + "..." + value[len(value)-4:]
			} else {
				displayValue = "****"
			}
		}

		fmt.Printf(SuccessStyle.Render("Set %s = %s\n"), key, displayValue)
		return nil
	},
}

// configGetCmd represents the config get subcommand
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.GetString(key)

		if value == "" {
			fmt.Printf(DimStyle.Render("%s is not set\n"), key)
			return nil
		}

		// Mask API keys in output
		displayValue := value
		if strings.Contains(strings.ToLower(key), "api_key") {
			if len(value) > 8 {
				displayValue = value[:4] + "..." + value[len(value)-4:]
			} else {
				displayValue = "****"
			}
		}

		fmt.Printf("%s = %s\n", key, displayValue)
		return nil
	},
}

// configListCmd represents the config list subcommand
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := viper.AllSettings()

		if len(settings) == 0 {
			fmt.Println(DimStyle.Render("No configuration found"))
			return nil
		}

		fmt.Println(TitleStyle.Render("Configuration:"))
		fmt.Println()

		printSettings("", settings)
		return nil
	},
}

// configPathCmd shows the config file path
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath := filepath.Join(home, ".gitcraft.yaml")
		fmt.Println(configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configPathCmd)
}

// writeConfig writes the current configuration to the config file
func writeConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".gitcraft.yaml")
	return viper.WriteConfigAs(configPath)
}

// printSettings recursively prints settings with proper indentation
func printSettings(prefix string, settings map[string]interface{}) {
	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			printSettings(fullKey, v)
		default:
			displayValue := fmt.Sprintf("%v", v)
			// Mask API keys
			if strings.Contains(strings.ToLower(fullKey), "api_key") {
				strVal := fmt.Sprintf("%v", v)
				if len(strVal) > 8 {
					displayValue = strVal[:4] + "..." + strVal[len(strVal)-4:]
				} else if strVal != "" {
					displayValue = "****"
				}
			}
			fmt.Printf("  %s = %s\n", InfoStyle.Render(fullKey), displayValue)
		}
	}
}
