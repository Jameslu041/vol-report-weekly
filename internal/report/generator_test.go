package report

import (
	"strings"
	"testing"

	"vol-report-weekly/internal/api"
	"vol-report-weekly/internal/storage"
)

func TestGenerator_Generate(t *testing.T) {
	gen := NewGenerator()

	current := &storage.Snapshot{
		Date: "2026-02-08",
		IVHistory: &api.IVHistoryResponse{
			Items: []api.IVHistoryData{
				{Tenor: "1D", IV: 45.5, HV: 42.0},
				{Tenor: "1W", IV: 48.0, HV: 44.0},
			},
		},
		IVRV: &api.IVRVResponse{
			Items: []api.IVRVData{
				{Period: "7D", IV: 50.0, RV: 45.0, VRP: 5.0},
			},
		},
		SkewChart: &api.SkewChartResponse{
			Items: []api.SkewData{
				{Tenor: "1M", Skew: 1.05},
			},
		},
		FIVMatrix: &api.FIVMatrixResponse{
			Items: []api.FIVMatrixEntry{
				{FromTenor: "1W", ToTenor: "1M", FIV: 48.5},
			},
		},
		OptionFlows: &api.OptionFlowsResponse{
			Data: api.OptionFlowsData{
				CallBuys:  1000,
				CallSells: 800,
				PutBuys:   600,
				PutSells:  500,
			},
		},
	}

	report, err := gen.Generate(current, nil)
	if err != nil {
		t.Fatalf("failed to generate report: %v", err)
	}

	if !strings.Contains(report, "BTC 波动率市场周报") {
		t.Error("report missing title")
	}
	if !strings.Contains(report, "2026-02-08") {
		t.Error("report missing date")
	}
	if !strings.Contains(report, "IV 期限结构概览") {
		t.Error("report missing IV section")
	}
	if !strings.Contains(report, "IV vs RV 分析") {
		t.Error("report missing IVRV section")
	}
	if !strings.Contains(report, "波动率偏斜") {
		t.Error("report missing skew section")
	}
	if !strings.Contains(report, "期权流向") {
		t.Error("report missing flows section")
	}
}

func TestGenerator_WithPrevious(t *testing.T) {
	gen := NewGenerator()

	current := &storage.Snapshot{
		Date: "2026-02-08",
		IVHistory: &api.IVHistoryResponse{
			Items: []api.IVHistoryData{
				{Tenor: "1D", IV: 50.0, HV: 42.0},
			},
		},
	}

	previous := &storage.Snapshot{
		Date: "2026-02-01",
		IVHistory: &api.IVHistoryResponse{
			Items: []api.IVHistoryData{
				{Tenor: "1D", IV: 45.0, HV: 40.0},
			},
		},
	}

	report, err := gen.Generate(current, previous)
	if err != nil {
		t.Fatalf("failed to generate report: %v", err)
	}

	if !strings.Contains(report, "上升") {
		t.Error("report should show IV increase")
	}
}

func TestAnalyzeIVChange(t *testing.T) {
	tests := []struct {
		current  float64
		previous float64
		want     string
	}{
		{50, 40, "大幅上升"},
		{40, 50, "大幅下降"},
		{50.5, 50, "基本持平"},
		{52, 50, "小幅上升"},
		{50, 0, "无历史数据"},
	}

	for _, tt := range tests {
		result := analyzeIVChange(tt.current, tt.previous)
		if !strings.Contains(result, tt.want) {
			t.Errorf("analyzeIVChange(%f, %f) = %s, want contains %s", tt.current, tt.previous, result, tt.want)
		}
	}
}

func TestAnalyzeVRP(t *testing.T) {
	tests := []struct {
		vrp  float64
		want string
	}{
		{10, "溢价偏高"},
		{3, "定价合理"},
		{-5, "折价"},
	}

	for _, tt := range tests {
		result := analyzeVRP(tt.vrp)
		if !strings.Contains(result, tt.want) {
			t.Errorf("analyzeVRP(%f) = %s, want contains %s", tt.vrp, result, tt.want)
		}
	}
}

func TestAnalyzeSkew(t *testing.T) {
	tests := []struct {
		skew float64
		want string
	}{
		{1.1, "看跌"},
		{1.0, "中性"},
		{0.9, "看涨"},
	}

	for _, tt := range tests {
		result := analyzeSkew(tt.skew)
		if !strings.Contains(result, tt.want) {
			t.Errorf("analyzeSkew(%f) = %s, want contains %s", tt.skew, result, tt.want)
		}
	}
}

func TestAnalyzeFlow(t *testing.T) {
	tests := []struct {
		callBuys float64
		putBuys  float64
		want     string
	}{
		{200, 100, "看涨"},
		{100, 200, "看跌"},
		{100, 100, "均衡"},
	}

	for _, tt := range tests {
		result := analyzeFlow(tt.callBuys, tt.putBuys)
		if !strings.Contains(result, tt.want) {
			t.Errorf("analyzeFlow(%f, %f) = %s, want contains %s", tt.callBuys, tt.putBuys, result, tt.want)
		}
	}
}
