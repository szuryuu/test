package fonnte

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kasiraiai/backend/config"
)

// Client adalah HTTP client untuk Fonnte WhatsApp API.
type Client struct {
	cfg        *config.Config
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:     cfg,
		baseURL: cfg.FonnteBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SendMessage mengirim pesan WhatsApp ke nomor tujuan melalui Fonnte API.
func (c *Client) SendMessage(to, message string) error {
	data := url.Values{}
	data.Set("target", to)
	data.Set("message", message)
	data.Set("delay", "1") // hindari rate limit ringan

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/send", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("fonnte.SendMessage: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", c.cfg.FonnteAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fonnte.SendMessage: do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fonnte.SendMessage: status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
