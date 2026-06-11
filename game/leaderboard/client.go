package leaderboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var DefaultLeaderboardURL string

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClientFromEnv() *Client {
	baseURL := strings.TrimSpace(os.Getenv("SNAKE_LEADERBOARD_URL"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(DefaultLeaderboardURL)
	}
	if baseURL == "" {
		return nil
	}
	return NewClient(baseURL)
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && strings.TrimSpace(c.BaseURL) != ""
}

func (c *Client) SubmitScore(ctx context.Context, req SubmitScoreRequest) (SubmitScoreResponse, error) {
	var out SubmitScoreResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/scores", req, &out); err != nil {
		return SubmitScoreResponse{}, err
	}
	return out, nil
}

func (c *Client) FetchVersionInfo(ctx context.Context) (VersionInfo, error) {
	var out VersionInfo
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/version", nil, &out); err != nil {
		return VersionInfo{}, err
	}
	return out, nil
}

func (c *Client) FetchLeaderboard(ctx context.Context, mode string, limit int) (LeaderboardResponse, error) {
	query := url.Values{}
	query.Set("mode", mode)
	query.Set("limit", fmt.Sprintf("%d", limit))

	var out LeaderboardResponse
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/leaderboard?"+query.Encode(), nil, &out); err != nil {
		return LeaderboardResponse{}, err
	}
	return out, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, in any, out any) error {
	if c == nil || c.BaseURL == "" {
		return fmt.Errorf("leaderboard client is disabled")
	}

	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if in != nil {
		req.Header.Set("content-type", "application/json")
	}

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr ErrorResponse
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Error != "" {
			return fmt.Errorf("leaderboard error: %s", apiErr.Error)
		}
		return fmt.Errorf("leaderboard status: %d", resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
