package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"kasiraiai/backend/config"
)

// Client adalah HTTP client untuk DeepSeek API (OpenAI-compatible).
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	baseURL    string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.DeepSeekTimeoutSecs) * time.Second,
		},
		baseURL: cfg.DeepSeekBaseURL,
	}
}

// ChatRequest adalah payload untuk /v1/chat/completions.
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse adalah response dari DeepSeek API.
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// Chat mengirim request ke DeepSeek API dan mengembalikan konten respons.
func (c *Client) Chat(systemPrompt, userPrompt string) (string, error) {
	reqBody := ChatRequest{
		Model: c.cfg.DeepSeekModel,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   c.cfg.DeepSeekMaxTokens,
		Temperature: 0.1, // rendah untuk parsing terstruktur
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ai.Client.Chat: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("ai.Client.Chat: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.DeepSeekAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ai.Client.Chat: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ai.Client.Chat: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ai.Client.Chat: status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("ai.Client.Chat: unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("ai.Client.Chat: empty choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
