package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClient_GetIVHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/deribit/datalab/iv_history" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("User-Agent") != "greeks_news_bot" {
			t.Errorf("unexpected user agent: %s", r.Header.Get("User-Agent"))
		}

		resp := map[string]any{
			"code":    200,
			"message": "",
			"data": []map[string]any{
				{
					"one_day":  45.5,
					"one_week": 48.0,
					"month1":   50.0,
					"month2":   52.0,
					"month3":   54.0,
					"month6":   56.0,
					"one_year": 58.0,
					"data": []map[string]any{
						{"day": 1, "hv": 42.0},
						{"day": 7, "hv": 44.0},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetIVHistory(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 7 {
		t.Errorf("expected 7 items, got %d", len(result.Items))
	}
	if result.Items[0].IV != 45.5 {
		t.Errorf("expected IV 45.5, got %f", result.Items[0].IV)
	}
}

func TestClient_GetIVRV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"code": 200,
			"data": map[string]any{
				"data": []map[string]any{
					{
						"datetime": 1234567890,
						"data": []map[string]any{
							{"day": 7, "iv": 50.0, "rv": 45.0, "vrp": 5.0},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetIVRV(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) == 0 {
		t.Error("expected at least one item")
	}
}

func TestClient_GetSkewChart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"code": 200,
			"data": []map[string]any{
				{
					"day1":    1.05,
					"days7":   1.02,
					"days30":  1.00,
					"days60":  0.98,
					"days90":  0.97,
					"days180": 0.96,
					"days365": 0.95,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetSkewChart(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 7 {
		t.Errorf("expected 7 items, got %d", len(result.Items))
	}
}

func TestClient_GetFIVMatrix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"code": 200,
			"data": map[string]any{
				"expiries": []map[string]any{
					{"label": "1W", "day": 7},
					{"label": "1M", "day": 30},
				},
				"matrix": [][]any{
					{nil, map[string]any{"Value": 48.5, "Diff": 1.0}},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetFIVMatrix(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}
}

func TestClient_GetOptionFlows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"code": 200,
			"data": map[string]any{
				"past24h": map[string]any{
					"call_buys":          1000.0,
					"call_sells":         800.0,
					"put_buys":           600.0,
					"put_sells":          500.0,
					"call_blocked_buys":  100.0,
					"call_blocked_sells": 50.0,
					"put_blocked_buys":   80.0,
					"put_blocked_sells":  40.0,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetOptionFlows(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data.CallBuys != 1000.0 {
		t.Errorf("expected CallBuys 1000.0, got %f", result.Data.CallBuys)
	}
}

func TestClient_RetryOnFailure(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := map[string]any{
			"code": 200,
			"data": []map[string]any{
				{
					"one_day":  45.5,
					"one_week": 48.0,
					"month1":   50.0,
					"month2":   52.0,
					"month3":   54.0,
					"month6":   56.0,
					"one_year": 58.0,
					"data":     []map[string]any{},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetIVHistory(context.Background())
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"code":    500,
			"message": "internal error",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetIVHistory(context.Background())
	if err == nil {
		t.Fatal("expected error for API error response")
	}
}
