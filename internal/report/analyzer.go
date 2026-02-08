package report

import "fmt"

func analyzeIVChange(current, previous float64) string {
	if previous == 0 {
		return "无历史数据对比"
	}
	change := ((current - previous) / previous) * 100
	abs := change
	if abs < 0 {
		abs = -abs
	}

	var trend string
	if change > 0 {
		trend = "上升"
	} else {
		trend = "下降"
	}

	if abs > 5 {
		return fmt.Sprintf("大幅%s %.1f%%", trend, abs)
	} else if abs >= 2 {
		return fmt.Sprintf("小幅%s %.1f%%", trend, abs)
	}
	return "基本持平"
}

func analyzeVRP(vrp float64) string {
	if vrp > 5 {
		return "期权溢价偏高，卖方有利"
	} else if vrp >= 0 {
		return "期权定价合理"
	}
	return "期权折价，买方有利"
}

func analyzeSkew(skew float64) string {
	if skew > 1.05 {
		return "明显看跌偏斜，市场偏谨慎"
	} else if skew >= 0.95 {
		return "偏斜中性"
	}
	return "看涨偏斜，市场偏乐观"
}

func analyzeFlow(callBuys, putBuys float64) string {
	if putBuys == 0 {
		if callBuys > 0 {
			return "看涨情绪浓厚"
		}
		return "无明显成交"
	}
	ratio := callBuys / putBuys
	if ratio > 1.5 {
		return "看涨情绪浓厚"
	} else if ratio < 0.67 {
		return "看跌情绪浓厚"
	}
	return "多空博弈均衡"
}
