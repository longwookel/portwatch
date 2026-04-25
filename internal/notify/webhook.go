package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookSender delivers events as JSON POST requests to a URL.
type WebhookSender struct {
	url    string
	client *http.Client
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Level   string `json:"level"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NewWebhookSender returns a WebhookSender that posts to url.
// A default 5-second timeout is applied when client is nil.
func NewWebhookSender(url string, client *http.Client) *WebhookSender {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &WebhookSender{url: url, client: client}
}

// Send marshals the event to JSON and POSTs it to the configured URL.
// Non-2xx responses are treated as errors.
func (ws *WebhookSender) Send(e Event) error {
	payload := webhookPayload{
		Level:   e.Level.String(),
		Title:   e.Title,
		Message: e.Message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := ws.client.Post(ws.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
