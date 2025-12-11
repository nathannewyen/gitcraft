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
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// OllamaProvider implements the Provider interface for Ollama (local models)
type OllamaProvider struct {
	baseURL string
	model   string
}

// Ollama API structures
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider() *OllamaProvider {
	baseURL := viper.GetString("ollama.url")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := viper.GetString("ollama.model")
	if model == "" {
		model = "llama3.2"
	}

	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
	}
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// IsConfigured returns true if Ollama is configured
// Ollama doesn't require an API key, just a running server
func (p *OllamaProvider) IsConfigured() bool {
	return p.baseURL != ""
}

// Generate generates a commit message using Ollama
func (p *OllamaProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	prompt := BuildPrompt(req)

	requestBody := ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal request: %v", ErrGenerationFailed, err)
	}

	// Security: Validate baseURL before constructing request URL
	parsedURL, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid base URL: %v", ErrGenerationFailed, err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("%w: base URL must use http or https scheme", ErrGenerationFailed)
	}

	requestURL := fmt.Sprintf("%s/api/generate", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrGenerationFailed, err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Ollama can be slow for local models, use longer timeout
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed (is Ollama running?): %v", ErrGenerationFailed, err)
	}
	defer resp.Body.Close()

	// Check HTTP status code before parsing response
	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: API error (status %d), failed to read response: %v", ErrGenerationFailed, resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("%w: API error (status %d): %s", ErrGenerationFailed, resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", ErrGenerationFailed, err)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse response: %v", ErrGenerationFailed, err)
	}

	if ollamaResp.Error != "" {
		return nil, fmt.Errorf("%w: %s", ErrGenerationFailed, ollamaResp.Error)
	}

	if ollamaResp.Response == "" {
		return nil, fmt.Errorf("%w: no response from Ollama", ErrGenerationFailed)
	}

	return ParseResponse(ollamaResp.Response), nil
}
