package analytics

type SummaryResponse struct {
	Sessions  int64 `json:"sessions"`
	PageViews int64 `json:"pageViews`
}

type Browser struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type BrowserListResponse struct {
	Items []*Browser `json:"items"`
}
