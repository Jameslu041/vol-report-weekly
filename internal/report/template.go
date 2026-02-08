package report

const reportTemplate = `# BTC 波动率市场周报

> 报告日期：{{.Date}}

## 一、IV 期限结构概览

| 期限 | 当前 IV | 历史 HV | 变化 |
|------|---------|---------|------|
{{range .IVHistory}}- {{.Tenor}} | {{printf "%.1f%%" .IV}} | {{printf "%.1f%%" .HV}} | {{.Change}}
{{end}}

## 二、IV vs RV 分析

| 周期 | IV | RV | VRP | 解读 |
|------|----|----|-----|------|
{{range .IVRV}}- {{.Period}} | {{printf "%.1f%%" .IV}} | {{printf "%.1f%%" .RV}} | {{printf "%.1f" .VRP}} | {{.Analysis}}
{{end}}

## 三、波动率偏斜（Skew）

| 期限 | Skew | 市场情绪 |
|------|------|----------|
{{range .Skew}}- {{.Tenor}} | {{printf "%.2f" .Skew}} | {{.Analysis}}
{{end}}

## 四、远期波动率矩阵（Forward IV）

| 起始期限 | 目标期限 | FIV |
|----------|----------|-----|
{{range .FIVMatrix}}- {{.FromTenor}} | {{.ToTenor}} | {{printf "%.1f%%" .FIV}}
{{end}}

## 五、期权流向（24h）

- Call 买入量：{{printf "%.0f" .OptionFlows.CallBuys}}
- Call 卖出量：{{printf "%.0f" .OptionFlows.CallSells}}
- Put 买入量：{{printf "%.0f" .OptionFlows.PutBuys}}
- Put 卖出量：{{printf "%.0f" .OptionFlows.PutSells}}
- 大宗买入：{{printf "%.0f" .OptionFlows.BlockBuys}}
- 大宗卖出：{{printf "%.0f" .OptionFlows.BlockSells}}

**市场情绪**：{{.FlowAnalysis}}

## 六、总结与交易观点

{{.Summary}}
`
