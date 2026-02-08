package api

// APIResponse is the common response wrapper
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// IVHistoryData represents IV history for different tenors (processed)
type IVHistoryData struct {
	Tenor string
	IV    float64
	HV    float64
}

// IVHistoryResponse contains the IV history data
type IVHistoryResponse struct {
	Items []IVHistoryData
}

// IVRVData represents IV vs RV comparison (processed)
type IVRVData struct {
	Period string
	IV     float64
	RV     float64
	VRP    float64
}

// IVRVResponse contains IV/RV data
type IVRVResponse struct {
	Items []IVRVData
}

// SkewData represents skew for a tenor (processed)
type SkewData struct {
	Tenor string
	Skew  float64
}

// SkewChartResponse contains skew chart data
type SkewChartResponse struct {
	Items []SkewData
}

// FIVMatrixEntry represents a forward IV entry
type FIVMatrixEntry struct {
	FromTenor string
	ToTenor   string
	FIV       float64
}

// FIVMatrixResponse contains forward IV matrix
type FIVMatrixResponse struct {
	Items []FIVMatrixEntry
}

// OptionFlowsData represents option flow data (processed)
type OptionFlowsData struct {
	CallBuys   float64
	CallSells  float64
	PutBuys    float64
	PutSells   float64
	BlockBuys  float64
	BlockSells float64
}

// OptionFlowsResponse contains option flows
type OptionFlowsResponse struct {
	Data OptionFlowsData
}
