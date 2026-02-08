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
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
	retryInterval  = 2 * time.Second
	userAgent      = "greeks_news_bot"
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

func (c *Client) doRequest(ctx context.Context, endpoint string, body any) (json.RawMessage, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			slog.Info("retrying request", "endpoint", endpoint, "attempt", attempt+1)
			time.Sleep(retryInterval)
		}

		data, err := c.doSingleRequest(ctx, endpoint, body)
		if err == nil {
			return data, nil
		}
		lastErr = err
		slog.Warn("request failed", "endpoint", endpoint, "attempt", attempt+1, "error", err)
	}

	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}

func (c *Client) doSingleRequest(ctx context.Context, endpoint string, body any) (json.RawMessage, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var wrapper struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if wrapper.Code != 200 {
		return nil, fmt.Errorf("api error: code=%d, message=%s", wrapper.Code, wrapper.Message)
	}

	return wrapper.Data, nil
}

func (c *Client) GetIVHistory(ctx context.Context) (*IVHistoryResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "7D",
		"exchange": "deribit",
	}

	data, err := c.doRequest(ctx, "/api/v1/deribit/datalab/iv_history", body)
	if err != nil {
		return nil, err
	}

	// Parse as array of entries
	var entries []struct {
		OneDay  float64 `json:"one_day"`
		OneWeek float64 `json:"one_week"`
		Month1  float64 `json:"month1"`
		Month2  float64 `json:"month2"`
		Month3  float64 `json:"month3"`
		Month6  float64 `json:"month6"`
		OneYear float64 `json:"one_year"`
		HVData  []struct {
			Day int     `json:"day"`
			HV  float64 `json:"hv"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse IV history: %w", err)
	}

	result := &IVHistoryResponse{}
	if len(entries) > 0 {
		latest := entries[len(entries)-1]
		hvMap := make(map[int]float64)
		for _, hv := range latest.HVData {
			hvMap[hv.Day] = hv.HV
		}

		result.Items = []IVHistoryData{
			{Tenor: "1D", IV: latest.OneDay, HV: hvMap[1]},
			{Tenor: "1W", IV: latest.OneWeek, HV: hvMap[7]},
			{Tenor: "1M", IV: latest.Month1, HV: hvMap[30]},
			{Tenor: "2M", IV: latest.Month2, HV: hvMap[60]},
			{Tenor: "3M", IV: latest.Month3, HV: hvMap[90]},
			{Tenor: "6M", IV: latest.Month6, HV: hvMap[180]},
			{Tenor: "1Y", IV: latest.OneYear, HV: hvMap[365]},
		}
	}

	return result, nil
}

func (c *Client) GetIVRV(ctx context.Context) (*IVRVResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "1M",
		"exchange": "deribit",
	}

	data, err := c.doRequest(ctx, "/api/v1/deribit/datalab/iv_rv", body)
	if err != nil {
		return nil, err
	}

	// Parse nested structure
	var wrapper struct {
		Data []struct {
			Datetime int64 `json:"datetime"`
			Data     []struct {
				Day int     `json:"day"`
				IV  float64 `json:"iv"`
				RV  float64 `json:"rv"`
				VRP float64 `json:"vrp"`
			} `json:"data"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse IV/RV: %w", err)
	}

	result := &IVRVResponse{}
	if len(wrapper.Data) > 0 {
		latest := wrapper.Data[len(wrapper.Data)-1]
		dayLabels := map[int]string{1: "1D", 3: "3D", 7: "1W", 15: "2W", 30: "1M", 90: "3M"}
		for _, item := range latest.Data {
			label := dayLabels[item.Day]
			if label == "" {
				label = fmt.Sprintf("%dD", item.Day)
			}
			result.Items = append(result.Items, IVRVData{
				Period: label,
				IV:     item.IV,
				RV:     item.RV,
				VRP:    item.IV - item.RV, // Calculate VRP
			})
		}
	}

	return result, nil
}

func (c *Client) GetSkewChart(ctx context.Context) (*SkewChartResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"gap":      "1M",
		"exchange": "deribit",
	}

	data, err := c.doRequest(ctx, "/api/v1/deribit/datalab/skew_chart", body)
	if err != nil {
		return nil, err
	}

	var entries []struct {
		Day1    float64 `json:"day1"`
		Days7   float64 `json:"days7"`
		Days30  float64 `json:"days30"`
		Days60  float64 `json:"days60"`
		Days90  float64 `json:"days90"`
		Days180 float64 `json:"days180"`
		Days365 float64 `json:"days365"`
	}

	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse skew: %w", err)
	}

	result := &SkewChartResponse{}
	if len(entries) > 0 {
		latest := entries[len(entries)-1]
		result.Items = []SkewData{
			{Tenor: "1D", Skew: latest.Day1},
			{Tenor: "1W", Skew: latest.Days7},
			{Tenor: "1M", Skew: latest.Days30},
			{Tenor: "2M", Skew: latest.Days60},
			{Tenor: "3M", Skew: latest.Days90},
			{Tenor: "6M", Skew: latest.Days180},
			{Tenor: "1Y", Skew: latest.Days365},
		}
	}

	return result, nil
}

func (c *Client) GetFIVMatrix(ctx context.Context) (*FIVMatrixResponse, error) {
	body := map[string]any{
		"currency": "BTC",
		"page_num": 0,
		"exchange": "deribit",
	}

	data, err := c.doRequest(ctx, "/api/v1/deribit/datalab/fiv_matrix", body)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Expiries []struct {
			Label string `json:"label"`
			Day   int    `json:"day"`
		} `json:"expiries"`
		Matrix [][]any `json:"matrix"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse FIV matrix: %w", err)
	}

	result := &FIVMatrixResponse{}
	// Extract diagonal entries (forward IVs)
	for i, row := range wrapper.Matrix {
		if i >= len(wrapper.Expiries) || i+1 >= len(row) {
			continue
		}
		cell := row[i+1]
		if cell == nil {
			continue
		}

		cellMap, ok := cell.(map[string]any)
		if !ok {
			continue
		}
		value, ok := cellMap["Value"].(float64)
		if !ok {
			continue
		}

		fromLabel := wrapper.Expiries[i].Label
		toLabel := ""
		if i+1 < len(wrapper.Expiries) {
			toLabel = wrapper.Expiries[i+1].Label
		}

		result.Items = append(result.Items, FIVMatrixEntry{
			FromTenor: fromLabel,
			ToTenor:   toLabel,
			FIV:       value,
		})
	}

	return result, nil
}

func (c *Client) GetOptionFlows(ctx context.Context) (*OptionFlowsResponse, error) {
	body := map[string]string{
		"currency": "BTC",
		"exchange": "deribit",
	}

	data, err := c.doRequest(ctx, "/api/v1/deribit/datalab/option_flows", body)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Past24h struct {
			CallBuys          float64 `json:"call_buys"`
			CallSells         float64 `json:"call_sells"`
			PutBuys           float64 `json:"put_buys"`
			PutSells          float64 `json:"put_sells"`
			CallBlockedBuys   float64 `json:"call_blocked_buys"`
			CallBlockedSells  float64 `json:"call_blocked_sells"`
			PutBlockedBuys    float64 `json:"put_blocked_buys"`
			PutBlockedSells   float64 `json:"put_blocked_sells"`
		} `json:"past24h"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse option flows: %w", err)
	}

	result := &OptionFlowsResponse{
		Data: OptionFlowsData{
			CallBuys:   wrapper.Past24h.CallBuys,
			CallSells:  wrapper.Past24h.CallSells,
			PutBuys:    wrapper.Past24h.PutBuys,
			PutSells:   wrapper.Past24h.PutSells,
			BlockBuys:  wrapper.Past24h.CallBlockedBuys + wrapper.Past24h.PutBlockedBuys,
			BlockSells: wrapper.Past24h.CallBlockedSells + wrapper.Past24h.PutBlockedSells,
		},
	}

	return result, nil
}
