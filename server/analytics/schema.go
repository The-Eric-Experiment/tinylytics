package analytics

type SummaryResponse struct {
	Sessions  int64 `json:"sessions"`
	PageViews int64 `json:"pageViews"`
}

type AnalyticsItem struct {
	Name  string  `json:"name"`
	Count int64   `json:"count"`
	Major *string `json:"major"`
	Minor *string `json:"minor"`
	Patch *string `json:"patch"`
}

type AnalyticsListResponse struct {
	Items []*AnalyticsItem `json:"items"`
}
