// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// ClaudeProvider implements the Provider interface for Anthropic Claude
type ClaudeProvider struct {
	apiKey string
	model  string
}

// Claude API structures
type claudeRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	Messages  []claudeMessage  `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
	Error   *claudeError    `json:"error,omitempty"`
}

type claudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type claudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider() *ClaudeProvider {
	apiKey := viper.GetString("claude.api_key")
	model := viper.GetString("claude.model")
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	return &ClaudeProvider{
		apiKey: apiKey,
		model:  model,
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// IsConfigured returns true if the provider has an API key configured
func (p *ClaudeProvider) IsConfigured() bool {
	return p.apiKey != ""
}

// Generate generates a commit message using Claude
func (p *ClaudeProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if !p.IsConfigured() {
		return nil, ErrNoAPIKey
	}

	prompt := BuildPrompt(req)

	requestBody := claudeRequest{
		Model:     p.model,
		MaxTokens: 500,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal request: %v", ErrGenerationFailed, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrGenerationFailed, err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %v", ErrGenerationFailed, err)
	}
	defer resp.Body.Close()

	// Check HTTP status code before parsing response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: API error (status %d): %s", ErrGenerationFailed, resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", ErrGenerationFailed, err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse response: %v", ErrGenerationFailed, err)
	}

	if claudeResp.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrGenerationFailed, claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("%w: no response from Claude", ErrGenerationFailed)
	}

	return ParseResponse(claudeResp.Content[0].Text), nil
}
