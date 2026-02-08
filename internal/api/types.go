package api

// APIResponse is the common response wrapper
type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// IVHistoryData represents IV history for different tenors
type IVHistoryData struct {
	Tenor string  `json:"tenor"`
	IV    float64 `json:"iv"`
	HV    float64 `json:"hv"`
	Date  string  `json:"date"`
}

// IVHistoryResponse contains the IV history data
type IVHistoryResponse struct {
	Items []IVHistoryData `json:"items"`
}

// IVRVData represents IV vs RV comparison
type IVRVData struct {
	Period string  `json:"period"`
	IV     float64 `json:"iv"`
	RV     float64 `json:"rv"`
	VRP    float64 `json:"vrp"`
}

// IVRVResponse contains IV/RV data
type IVRVResponse struct {
	Items []IVRVData `json:"items"`
}

// SkewData represents skew for a tenor
type SkewData struct {
	Tenor string  `json:"tenor"`
	Skew  float64 `json:"skew"`
	Date  string  `json:"date"`
}

// SkewChartResponse contains skew chart data
type SkewChartResponse struct {
	Items []SkewData `json:"items"`
}

// FIVMatrixEntry represents a forward IV entry
type FIVMatrixEntry struct {
	FromTenor string  `json:"from_tenor"`
	ToTenor   string  `json:"to_tenor"`
	FIV       float64 `json:"fiv"`
}

// FIVMatrixResponse contains forward IV matrix
type FIVMatrixResponse struct {
	Items []FIVMatrixEntry `json:"items"`
}

// OptionFlowsData represents option flow data
type OptionFlowsData struct {
	CallBuys   float64 `json:"call_buys"`
	CallSells  float64 `json:"call_sells"`
	PutBuys    float64 `json:"put_buys"`
	PutSells   float64 `json:"put_sells"`
	BlockBuys  float64 `json:"block_buys"`
	BlockSells float64 `json:"block_sells"`
}

// OptionFlowsResponse contains option flows
type OptionFlowsResponse struct {
	Data OptionFlowsData `json:"data"`
}
