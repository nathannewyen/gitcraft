// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"
)

func TestMaskAPIKey_LongKey(t *testing.T) {
	key := "sk-1234567890abcdefghijklmnop"

	result := maskAPIKey(key)

	// Should show first 4 and last 4 characters
	if result != "sk-1...mnop" {
		t.Errorf("Expected 'sk-1...mnop', got '%s'", result)
	}
}

func TestMaskAPIKey_ShortKey(t *testing.T) {
	key := "short"

	result := maskAPIKey(key)

	// Short keys should be fully masked
	if result != "****" {
		t.Errorf("Expected '****', got '%s'", result)
	}
}

func TestMaskAPIKey_ExactlyEightChars(t *testing.T) {
	key := "12345678"

	result := maskAPIKey(key)

	// 8 char keys should be fully masked (not > 8)
	if result != "****" {
		t.Errorf("Expected '****', got '%s'", result)
	}
}

func TestMaskAPIKey_NineChars(t *testing.T) {
	key := "123456789"

	result := maskAPIKey(key)

	// 9 char keys should show partial
	if result != "1234...6789" {
		t.Errorf("Expected '1234...6789', got '%s'", result)
	}
}

func TestMaskAPIKey_EmptyKey(t *testing.T) {
	result := maskAPIKey("")

	if result != "" {
		t.Errorf("Expected empty string for empty key, got '%s'", result)
	}
}

func TestMaskAPIKey_VeryLongKey(t *testing.T) {
	key := "sk-proj-abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJ"

	result := maskAPIKey(key)

	// Should still work for very long keys
	if len(result) > 15 {
		t.Errorf("Masked key should be short, got '%s' (len=%d)", result, len(result))
	}
	// Last 4 chars of key are "GHIJ"
	if result != "sk-p...GHIJ" {
		t.Errorf("Expected 'sk-p...GHIJ', got '%s'", result)
	}
}

func TestStyles_Defined(t *testing.T) {
	// Test that style constants are defined and don't panic
	_ = TitleStyle.Render("test")
	_ = SuccessStyle.Render("test")
	_ = ErrorStyle.Render("test")
	_ = InfoStyle.Render("test")
	_ = DimStyle.Render("test")
}
