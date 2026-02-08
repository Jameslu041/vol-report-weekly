package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultTimeout    = 30 * time.Second
	maxRetries        = 3
	retryInterval     = 2 * time.Second
	userAgent         = "greeks_news_bot"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			slog.Info("retrying request", "endpoint", endpoint, "attempt", attempt+1)
			time.Sleep(retryInterval)
		}

		err := c.doSingleRequest(ctx, endpoint, body, result)
		if err == nil {
			return nil
		}
		lastErr = err
		slog.Warn("request failed", "endpoint", endpoint, "attempt", attempt+1, "error", err)
	}

	return fmt.Errorf("all retries failed: %w", lastErr)
}

func (c *Client) doSingleRequest(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

func (c *Client) GetIVHistory(ctx context.Context) (*IVHistoryResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "7D",
		"exchange": "deribit",
	}

	var resp APIResponse[IVHistoryResponse]
	if err := c.doRequest(ctx, "/api/v1/deribit/datalab/iv_history", body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

func (c *Client) GetIVRV(ctx context.Context) (*IVRVResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "1M",
		"exchange": "deribit",
	}

	var resp APIResponse[IVRVResponse]
	if err := c.doRequest(ctx, "/api/v1/deribit/datalab/iv_rv", body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

func (c *Client) GetSkewChart(ctx context.Context) (*SkewChartResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "1M",
		"exchange": "deribit",
	}

	var resp APIResponse[SkewChartResponse]
	if err := c.doRequest(ctx, "/api/v1/deribit/datalab/skew_chart", body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

func (c *Client) GetFIVMatrix(ctx context.Context) (*FIVMatrixResponse, error) {
	body := map[string]interface{}{
		"currency": "BTC",
		"page_num": 0,
		"exchange": "deribit",
	}

	var resp APIResponse[FIVMatrixResponse]
	if err := c.doRequest(ctx, "/api/v1/deribit/datalab/fiv_matrix", body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

func (c *Client) GetOptionFlows(ctx context.Context) (*OptionFlowsResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"exchange": "deribit",
	}

	var resp APIResponse[OptionFlowsResponse]
	if err := c.doRequest(ctx, "/api/v1/deribit/datalab/option_flows", body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}
