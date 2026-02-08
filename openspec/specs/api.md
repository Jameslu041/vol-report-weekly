# Spec: Deribit DataLab API Client

## Overview

Deribit DataLab API 客户端，用于获取 BTC 期权波动率数据。

## API Endpoints

Base URL: `https://dev-api.yazhan.vip`
Header: `User-Agent: greeks_news_bot`

### 1. IV History

```
POST /api/v1/deribit/datalab/iv_history
Body: {"currency": "BTC", "gap": "7D", "exchange": "deribit"}
```

Response: 各期限（1D/1W/1M/2M/3M/6M/1Y）的 IV 数据及 HV 对照

### 2. IV/RV

```
POST /api/v1/deribit/datalab/iv_rv
Body: {"currency": "BTC", "gap": "1M", "exchange": "deribit"}
```

Response: 各周期的 IV、RV、VRP

### 3. Skew Chart

```
POST /api/v1/deribit/datalab/skew_chart
Body: {"currency": "BTC", "gap": "1M", "exchange": "deribit"}
```

Response: 各期限偏斜值

### 4. FIV Matrix

```
POST /api/v1/deribit/datalab/fiv_matrix
Body: {"currency": "BTC", "page_num": 0, "exchange": "deribit"}
```

Response: 远期 IV 期限结构矩阵

### 5. Option Flows

```
POST /api/v1/deribit/datalab/option_flows
Body: {"currency": "BTC", "exchange": "deribit"}
```

Response: 24h 内 call/put 买卖量

## Client Interface

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    userAgent  string
}

func NewClient(baseURL string) *Client
func (c *Client) GetIVHistory(ctx context.Context) (*IVHistoryResponse, error)
func (c *Client) GetIVRV(ctx context.Context) (*IVRVResponse, error)
func (c *Client) GetSkewChart(ctx context.Context) (*SkewChartResponse, error)
func (c *Client) GetFIVMatrix(ctx context.Context) (*FIVMatrixResponse, error)
func (c *Client) GetOptionFlows(ctx context.Context) (*OptionFlowsResponse, error)
```

## Retry Policy

- Max retries: 3
- Retry interval: 2 seconds
- Timeout: 30 seconds per request
