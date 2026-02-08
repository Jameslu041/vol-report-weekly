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

		resp := APIResponse[IVHistoryResponse]{
			Code:    200,
			Message: "",
			Data: IVHistoryResponse{
				Items: []IVHistoryData{
					{Tenor: "1D", IV: 45.5, HV: 42.0, Date: "2026-02-08"},
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
	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].IV != 45.5 {
		t.Errorf("expected IV 45.5, got %f", result.Items[0].IV)
	}
}

func TestClient_GetIVRV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[IVRVResponse]{
			Code: 200,
			Data: IVRVResponse{
				Items: []IVRVData{
					{Period: "7D", IV: 50.0, RV: 45.0, VRP: 5.0},
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
	if result.Items[0].VRP != 5.0 {
		t.Errorf("expected VRP 5.0, got %f", result.Items[0].VRP)
	}
}

func TestClient_GetSkewChart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[SkewChartResponse]{
			Code: 200,
			Data: SkewChartResponse{
				Items: []SkewData{
					{Tenor: "1M", Skew: 1.05, Date: "2026-02-08"},
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
	if result.Items[0].Skew != 1.05 {
		t.Errorf("expected skew 1.05, got %f", result.Items[0].Skew)
	}
}

func TestClient_GetFIVMatrix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[FIVMatrixResponse]{
			Code: 200,
			Data: FIVMatrixResponse{
				Items: []FIVMatrixEntry{
					{FromTenor: "1W", ToTenor: "1M", FIV: 48.5},
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
	if result.Items[0].FIV != 48.5 {
		t.Errorf("expected FIV 48.5, got %f", result.Items[0].FIV)
	}
}

func TestClient_GetOptionFlows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse[OptionFlowsResponse]{
			Code: 200,
			Data: OptionFlowsResponse{
				Data: OptionFlowsData{
					CallBuys:  1000.0,
					CallSells: 800.0,
					PutBuys:   600.0,
					PutSells:  500.0,
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
		resp := APIResponse[IVHistoryResponse]{
			Code: 200,
			Data: IVHistoryResponse{Items: []IVHistoryData{}},
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
		resp := APIResponse[IVHistoryResponse]{
			Code:    500,
			Message: "internal error",
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
