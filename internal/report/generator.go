package report

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"vol-report-weekly/internal/storage"
)

type IVHistoryItem struct {
	Tenor  string
	IV     float64
	HV     float64
	Change string
}

type IVRVItem struct {
	Period   string
	IV       float64
	RV       float64
	VRP      float64
	Analysis string
}

type SkewItem struct {
	Tenor    string
	Skew     float64
	Analysis string
}

type FIVItem struct {
	FromTenor string
	ToTenor   string
	FIV       float64
}

type OptionFlowsData struct {
	CallBuys   float64
	CallSells  float64
	PutBuys    float64
	PutSells   float64
	BlockBuys  float64
	BlockSells float64
}

type ReportData struct {
	Date         string
	IVHistory    []IVHistoryItem
	IVRV         []IVRVItem
	Skew         []SkewItem
	FIVMatrix    []FIVItem
	OptionFlows  OptionFlowsData
	FlowAnalysis string
	Summary      string
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(current, previous *storage.Snapshot) (string, error) {
	data := g.buildReportData(current, previous)

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

func (g *Generator) buildReportData(current, previous *storage.Snapshot) ReportData {
	data := ReportData{
		Date: current.Date,
	}

	// IV History
	if current.IVHistory != nil {
		for _, item := range current.IVHistory.Items {
			prevIV := 0.0
			if previous != nil && previous.IVHistory != nil {
				for _, p := range previous.IVHistory.Items {
					if p.Tenor == item.Tenor {
						prevIV = p.IV
						break
					}
				}
			}
			data.IVHistory = append(data.IVHistory, IVHistoryItem{
				Tenor:  item.Tenor,
				IV:     item.IV,
				HV:     item.HV,
				Change: analyzeIVChange(item.IV, prevIV),
			})
		}
	}

	// IV/RV
	if current.IVRV != nil {
		for _, item := range current.IVRV.Items {
			data.IVRV = append(data.IVRV, IVRVItem{
				Period:   item.Period,
				IV:       item.IV,
				RV:       item.RV,
				VRP:      item.VRP,
				Analysis: analyzeVRP(item.VRP),
			})
		}
	}

	// Skew
	if current.SkewChart != nil {
		for _, item := range current.SkewChart.Items {
			data.Skew = append(data.Skew, SkewItem{
				Tenor:    item.Tenor,
				Skew:     item.Skew,
				Analysis: analyzeSkew(item.Skew),
			})
		}
	}

	// FIV Matrix
	if current.FIVMatrix != nil {
		for _, item := range current.FIVMatrix.Items {
			data.FIVMatrix = append(data.FIVMatrix, FIVItem{
				FromTenor: item.FromTenor,
				ToTenor:   item.ToTenor,
				FIV:       item.FIV,
			})
		}
	}

	// Option Flows
	if current.OptionFlows != nil {
		data.OptionFlows = OptionFlowsData{
			CallBuys:   current.OptionFlows.Data.CallBuys,
			CallSells:  current.OptionFlows.Data.CallSells,
			PutBuys:    current.OptionFlows.Data.PutBuys,
			PutSells:   current.OptionFlows.Data.PutSells,
			BlockBuys:  current.OptionFlows.Data.BlockBuys,
			BlockSells: current.OptionFlows.Data.BlockSells,
		}
		data.FlowAnalysis = analyzeFlow(current.OptionFlows.Data.CallBuys, current.OptionFlows.Data.PutBuys)
	}

	// Summary
	data.Summary = g.generateSummary(data)

	return data
}

func (g *Generator) generateSummary(data ReportData) string {
	var points []string

	// IV trend
	if len(data.IVHistory) > 0 {
		if strings.Contains(data.IVHistory[0].Change, "上升") {
			points = append(points, "短期波动率上升")
		} else if strings.Contains(data.IVHistory[0].Change, "下降") {
			points = append(points, "短期波动率下降")
		}
	}

	// VRP
	if len(data.IVRV) > 0 {
		if data.IVRV[0].VRP > 5 {
			points = append(points, "期权溢价偏高，可关注卖方策略")
		} else if data.IVRV[0].VRP < 0 {
			points = append(points, "期权折价，可关注买方策略")
		}
	}

	// Flow
	points = append(points, "市场情绪："+data.FlowAnalysis)

	if len(points) == 0 {
		return "市场整体波动率处于正常水平，建议观望。"
	}

	return strings.Join(points, "；") + "。"
}
