// Copyright 2025 Nathan Nguyen
// SPDX-License-Identifier: MIT

package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	client *openai.Client
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider() *OpenAIProvider {
	apiKey := viper.GetString("openai.api_key")
	if apiKey == "" {
		return &OpenAIProvider{}
	}

	model := viper.GetString("openai.model")
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIProvider{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// IsConfigured returns true if the provider has an API key configured
func (p *OpenAIProvider) IsConfigured() bool {
	return p.client != nil
}

// Generate generates a commit message using OpenAI
func (p *OpenAIProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if !p.IsConfigured() {
		return nil, ErrNoAPIKey
	}

	prompt := BuildPrompt(req)

	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   500,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGenerationFailed, err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("%w: no response from OpenAI", ErrGenerationFailed)
	}

	return ParseResponse(resp.Choices[0].Message.Content), nil
}
