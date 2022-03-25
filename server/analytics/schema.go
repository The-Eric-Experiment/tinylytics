package analytics

type SummaryResponse struct {
	Sessions           int64   `json:"sessions"`
	PageViews          int64   `json:"pageViews"`
	AvgSessionDuration float64 `json:"avgSessionDuration"`
	BounceRate         int64   `json:"bounceRate"`
}

type AnalyticsItem struct {
	Value     string `json:"value"`
	Count     int64  `json:"count"`
	Drillable int64  `json:"drillable"`
}

type AnalyticsListResponse struct {
	PreviousFilters []string         `json:"previousFilters"`
	Items           []*AnalyticsItem `json:"items"`
}
