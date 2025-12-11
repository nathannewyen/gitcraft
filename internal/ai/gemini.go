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

	"github.com/spf13/viper"
)

// GeminiProvider implements the Provider interface for Google Gemini
type GeminiProvider struct {
	apiKey string
	model  string
}

// Gemini API structures
type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *geminiError      `json:"error,omitempty"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider() *GeminiProvider {
	apiKey := viper.GetString("gemini.api_key")
	model := viper.GetString("gemini.model")
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
	}
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// IsConfigured returns true if the provider has an API key configured
func (p *GeminiProvider) IsConfigured() bool {
	return p.apiKey != ""
}

// Generate generates a commit message using Gemini
func (p *GeminiProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if !p.IsConfigured() {
		return nil, ErrNoAPIKey
	}

	prompt := BuildPrompt(req)

	requestBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:     0.3,
			MaxOutputTokens: 500,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal request: %v", ErrGenerationFailed, err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrGenerationFailed, err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %v", ErrGenerationFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", ErrGenerationFailed, err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse response: %v", ErrGenerationFailed, err)
	}

	if geminiResp.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrGenerationFailed, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("%w: no response from Gemini", ErrGenerationFailed)
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("%w: empty response from Gemini", ErrGenerationFailed)
	}

	return ParseResponse(geminiResp.Candidates[0].Content.Parts[0].Text), nil
}
