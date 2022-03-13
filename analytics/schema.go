package analytics

type SummaryResponse struct {
	Sessions  int64 `json:"sessions"`
	PageViews int64 `json:"pageViews`
}

type Browser struct {
	Name         string `json:"name"`
	Count        int64  `json:"count"`
	BrowserMajor *int64 `json:"browserMajor"`
	BrowserMinor *int64 `json:"browserMinor"`
	BrowserPatch *int64 `json:"browserPatch"`
}

type BrowserListResponse struct {
	Items []*Browser `json:"items"`
}
