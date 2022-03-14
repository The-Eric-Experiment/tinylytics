package analytics

type SummaryResponse struct {
	Sessions  int64 `json:"sessions"`
	PageViews int64 `json:"pageViews`
}

type AnalyticsItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
	Major *int64 `json:"major"`
	Minor *int64 `json:"minor"`
	Patch *int64 `json:"patch"`
}

type AnalyticsListResponse struct {
	Items []*AnalyticsItem `json:"items"`
}
