# Spec: Data Storage

## Overview

数据快照存储模块，将 API 响应保存为 JSON 文件。

## Snapshot Structure

```go
type Snapshot struct {
    Date        string                   `json:"date"`
    IVHistory   *api.IVHistoryResponse   `json:"iv_history"`
    IVRV        *api.IVRVResponse        `json:"iv_rv"`
    SkewChart   *api.SkewChartResponse   `json:"skew_chart"`
    FIVMatrix   *api.FIVMatrixResponse   `json:"fiv_matrix"`
    OptionFlows *api.OptionFlowsResponse `json:"option_flows"`
}
```

## Storage Interface

```go
type Storage struct {
    dataDir string
}

func NewStorage(dataDir string) *Storage
func (s *Storage) Save(snapshot *Snapshot) error
func (s *Storage) Load(date string) (*Snapshot, error)
func (s *Storage) LoadPrevious(currentDate string) (*Snapshot, error)
```

## File Format

- Path: `{dataDir}/snapshots/YYYY-MM-DD.json`
- Encoding: UTF-8 JSON with indentation
