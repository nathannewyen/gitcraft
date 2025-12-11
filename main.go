// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package main

import (
	"os"

	"github.com/nathannewyen/gitcraft/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
