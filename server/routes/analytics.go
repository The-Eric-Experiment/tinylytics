package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"tinylytics/analytics"
	"tinylytics/config"
	"tinylytics/db"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
)

// Template data structures
type TemplateData struct {
	Title         string
	Domain        string
	Websites      []config.WebsiteConfig
	CurrentPeriod string
	Periods       []PeriodOption
	ActiveFilters []ActiveFilter
	Summary       *SummaryData
	QueryString   string
}

type SummaryData struct {
	Sessions           int64
	PageViews          int64
	AvgSessionDuration string
	BounceRate         int64
}

type PeriodOption struct {
	Value string
	Label string
}

type ActiveFilter struct {
	Name        string
	Value       string
	RemoveQuery string
}

type AnalyticsItemWithIcon struct {
	*analytics.AnalyticsItem `json:",inline"`
	Icon                     string `json:"icon,omitempty"`
	CountryName              string `json:"countryName,omitempty"`
	CountryCode              string `json:"countryCode,omitempty"`
	Label                    string `json:"label,omitempty"`
	IsClickable              bool   `json:"isClickable,omitempty"`
	FaviconURL               string `json:"faviconUrl,omitempty"`
	FormattedValue           string `json:"formattedValue,omitempty"`
	FilterKey                string `json:"filterKey,omitempty"`
	FilterValue              string `json:"filterValue,omitempty"`
}

func arrayFromRows(rows *sql.Rows) []*analytics.AnalyticsItem {
	if rows == nil {
		log.Printf("ERROR: rows is nil in arrayFromRows")
		return make([]*analytics.AnalyticsItem, 0)
	}
	defer rows.Close()

	list := make([]*analytics.AnalyticsItem, 0)

	for rows.Next() {
		var output analytics.AnalyticsItem

		// All analytics queries return: value, count, drillable
		err := rows.Scan(&output.Value, &output.Count, &output.Drillable)
		if err != nil {
			log.Printf("ERROR: Failed to scan row: %v", err)
			continue
		}

		list = append(list, &output)
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Error iterating rows: %v", err)
	}

	return list
}

func getDB(c *gin.Context) *db.Database {
	domain := c.Query("site")
	if domain == "" {
		domain, _ = c.Params.Get("domain")
	}

	if domain == "" {
		c.String(http.StatusBadRequest, "Domain required")
		return nil
	}

	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return nil
	}

	database := db.Database{}
	database.Connect(dbFile)
	return &database
}

// GetAnalyticsPage - main page handler
func GetAnalyticsPage(c *gin.Context) {
	domain := c.Query("site")
	if domain == "" {
		domain = getDefaultDomain()
		if domain != "" {
			c.Redirect(http.StatusFound, "/?site="+domain)
			return
		}
	}

	// Build template data WITHOUT fetching database data
	data := buildTemplateDataShell(c, domain)
	c.HTML(http.StatusOK, "analytics.html", data)
}

// GetSummaries - returns HTML template instead of JSON
func GetSummaries(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	sessions := database.GetSessions(c)
	pageViews := database.GetPageViews(c)
	avgSessionDuration := database.GetAvgSessionDuration(c)
	bounceRate := database.GetBounceRate(c)

	summary := &SummaryData{
		Sessions:           sessions,
		PageViews:          pageViews,
		AvgSessionDuration: formatDuration(avgSessionDuration),
		BounceRate:         bounceRate,
	}

	data := map[string]interface{}{
		"Summary":     summary,
		"QueryString": buildQueryString(c),
	}

	c.HTML(http.StatusOK, "summary.html", data)
}

// GetBrowsers - returns HTML template instead of JSON
func GetBrowsers(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	rows, err := database.GetBrowsers(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
		return
	}

	items := arrayFromRows(rows)

	previousFilters := make([]string, 0)
	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	if hasBrowser {
		previousFilters = append(previousFilters, browser)
	}
	if hasBrowserVersion {
		bver := strings.Split(browserVersion, "/")
		previousFilters = append(previousFilters, bver...)
	}

	itemsWithIcons := processBrowserItems(items, browser, browserVersion, previousFilters, len(items) > 1)

	data := map[string]interface{}{
		"Domain":          domain,
		"CurrentPeriod":   c.DefaultQuery("p", "24h"),
		"PreviousFilters": previousFilters,
		"Items":           itemsWithIcons,
		"QueryString":     buildQueryString(c),
		"FilterPrimary":   "b",
		"FilterSecondary": "bv",
	}

	c.HTML(http.StatusOK, "browsers-table.html", data)
}

// GetOSs - returns HTML template instead of JSON
func GetOSs(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	rows, err := database.GetOSs(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get OSs")
		return
	}

	items := arrayFromRows(rows)

	previousFilters := make([]string, 0)
	os, hasOs := c.GetQuery("os")
	osVersion, hasOsVersion := c.GetQuery("osv")

	if hasOs {
		previousFilters = append(previousFilters, os)
	}
	if hasOsVersion {
		osver := strings.Split(osVersion, "/")
		previousFilters = append(previousFilters, osver...)
	}

	processedItems := processOSItems(items, os, osVersion, previousFilters, len(items) > 1)

	data := map[string]interface{}{
		"Domain":          domain,
		"CurrentPeriod":   c.DefaultQuery("p", "24h"),
		"PreviousFilters": previousFilters,
		"Items":           processedItems,
		"QueryString":     buildQueryString(c),
		"FilterPrimary":   "os",
		"FilterSecondary": "osv",
	}

	c.HTML(http.StatusOK, "os-table.html", data)
}

// GetCountries - returns HTML template instead of JSON
func GetCountries(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	rows, err := database.GetCountries(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Countries")
		return
	}

	items := arrayFromRows(rows)

	country, hasCountry := c.GetQuery("c")
	previousFilters := make([]string, 0)
	if hasCountry {
		previousFilters = append(previousFilters, getCountryName(country))
	}

	// Process all items for the map
	processedItemsAll := processCountryItems(items, country, previousFilters, len(items) > 1)

	// Get first 20 items for the table
	tableItems := processedItemsAll
	if len(tableItems) > 20 {
		tableItems = tableItems[:20]
	}

	data := map[string]interface{}{
		"Domain":          domain,
		"CurrentPeriod":   c.DefaultQuery("p", "24h"),
		"PreviousFilters": previousFilters,
		"Items":           tableItems,
		"MapItems":        processedItemsAll,
		"QueryString":     buildQueryString(c),
		"FilterPrimary":   "c",
	}

	c.HTML(http.StatusOK, "countries-table.html", data)
}

// GetReferrers - returns HTML template instead of JSON
func GetReferrers(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	rows, err := database.GetReferrers(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Referrers")
		return
	}

	items := arrayFromRows(rows)

	referrer, hasReferrer := c.GetQuery("r")
	previousFilters := make([]string, 0)
	if hasReferrer {
		previousFilters = append(previousFilters, referrer)
	}

	referrerPath := c.Query("rfp")
	processedItems := processReferrerItems(items, referrer, referrerPath, previousFilters, len(items) > 1)

	data := map[string]interface{}{
		"Domain":          domain,
		"CurrentPeriod":   c.DefaultQuery("p", "24h"),
		"PreviousFilters": previousFilters,
		"Items":           processedItems,
		"QueryString":     buildQueryString(c),
		"FilterPrimary":   "r",
		"FilterSecondary": "rfp",
	}

	c.HTML(http.StatusOK, "referrers-table.html", data)
}

// GetPages - returns HTML template instead of JSON
func GetPages(c *gin.Context) {
	domain := c.Query("site")
	c.Params = append(c.Params, gin.Param{Key: "domain", Value: domain})

	database := getDB(c)
	if database == nil {
		return
	}
	defer database.Close()

	rows, err := database.GetPages(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Pages")
		return
	}

	items := arrayFromRows(rows)

	path, hasPath := c.GetQuery("pg")
	previousFilters := make([]string, 0)
	if hasPath {
		previousFilters = append(previousFilters, path)
	}

	processedItems := processPageItems(items, path, previousFilters, len(items) > 1)

	data := map[string]interface{}{
		"Domain":          domain,
		"CurrentPeriod":   c.DefaultQuery("p", "24h"),
		"PreviousFilters": previousFilters,
		"Items":           processedItems,
		"QueryString":     buildQueryString(c),
		"FilterPrimary":   "pg",
	}

	c.HTML(http.StatusOK, "pages-table.html", data)
}

func GetWebsites(c *gin.Context) {
	sites := config.Config.Websites
	c.IndentedJSON(http.StatusOK, &sites)
}

// Template helper functions (moved from templates.go)

func buildTemplateDataShell(c *gin.Context, domain string) *TemplateData {
	return &TemplateData{
		Title:         domain,
		Domain:        domain,
		Websites:      config.Config.Websites,
		CurrentPeriod: c.DefaultQuery("p", "24h"),
		Periods:       getPeriodOptions(),
		ActiveFilters: buildActiveFilters(c),
		QueryString:   buildQueryString(c),
	}
}

func buildQueryString(c *gin.Context) string {
	query := c.Request.URL.Query()
	query.Del("site")
	query.Del("p")
	if len(query) == 0 {
		return ""
	}
	return "&" + query.Encode()
}

var dependantFilters = map[string][]string{
	"b":  {"bv"},
	"os": {"osv"},
	"r":  {"rfp"},
}

var showAsSameFilter = [][]string{
	{"r", "rfp"},
}

func buildActiveFilters(c *gin.Context) []ActiveFilter {
	query := c.Request.URL.Query()

	// Normalize query params to avoid duplicates
	normalizedQuery := url.Values{}
	normalizedQuery.Set("site", c.DefaultQuery("site", ""))
	normalizedQuery.Set("p", c.DefaultQuery("p", "24h"))
	for key := range query {
		if key != "site" && key != "p" {
			normalizedQuery.Set(key, query.Get(key))
		}
	}
	query = normalizedQuery

	allFilters := map[string]string{
		"b":   query.Get("b"),
		"bv":  query.Get("bv"),
		"os":  query.Get("os"),
		"osv": query.Get("osv"),
		"c":   query.Get("c"),
		"r":   query.Get("r"),
		"rfp": query.Get("rfp"),
		"pg":  query.Get("pg"),
	}

	filterNames := map[string]string{
		"b":   "Browser",
		"bv":  "Browser Version",
		"os":  "OS",
		"osv": "OS Version",
		"c":   "Country",
		"r":   "Referrer",
		"rfp": "Referrer",
		"pg":  "Page",
	}

	presentKeys := []string{}
	for key, val := range allFilters {
		if val != "" {
			presentKeys = append(presentKeys, key)
		}
	}

	reducedKeys := []string{}
	for _, key := range presentKeys {
		showAsSame := findShowAsSame(key)
		if showAsSame == nil {
			reducedKeys = append(reducedKeys, key)
			continue
		}

		inAccIndex := -1
		for i, existing := range reducedKeys {
			if containsString(showAsSame, existing) {
				inAccIndex = i
				break
			}
		}

		if inAccIndex == -1 {
			reducedKeys = append(reducedKeys, key)
		} else {
			indexCurrent := indexOf(showAsSame, key)
			indexExisting := indexOf(showAsSame, reducedKeys[inAccIndex])
			if indexCurrent > indexExisting {
				reducedKeys[inAccIndex] = key
			}
		}
	}

	filters := []ActiveFilter{}
	for _, key := range reducedKeys {
		value := allFilters[key]
		if value == "" {
			continue
		}

		displayValue := value
		if key == "c" {
			displayValue = getCountryName(value)
		}

		q := cloneQuery(query)
		q.Del(key)
		if deps, ok := dependantFilters[key]; ok {
			for _, dep := range deps {
				q.Del(dep)
			}
		}

		filters = append(filters, ActiveFilter{
			Name:        filterNames[key],
			Value:       displayValue,
			RemoveQuery: q.Encode(),
		})
	}

	return filters
}

func findShowAsSame(key string) []string {
	for _, group := range showAsSameFilter {
		if containsString(group, key) {
			return group
		}
	}
	return nil
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func indexOf(slice []string, str string) int {
	for i, s := range slice {
		if s == str {
			return i
		}
	}
	return -1
}

func cloneQuery(q url.Values) url.Values {
	n := make(url.Values)
	for k, v := range q {
		n[k] = v
	}
	return n
}

func getPeriodOptions() []PeriodOption {
	return []PeriodOption{
		{"24h", "Last 24 Hours"},
		{"today", "Today"},
		{"yesterday", "Yesterday"},
		{"7d", "Last 7 Days"},
		{"week", "This Week"},
		{"lastweek", "Last Week"},
		{"30d", "Last 30 Days"},
		{"90d", "Last 90 Days"},
		{"month", "This Month"},
		{"lastmonth", "Last Month"},
		{"year", "This Year"},
		{"lastyear", "Last Year"},
		{"alltime", "All Time"},
	}
}

func getDefaultDomain() string {
	if len(config.Config.Websites) > 0 {
		return config.Config.Websites[0].Domain
	}
	return ""
}

func formatDuration(seconds float64) string {
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm %ds", minutes, secs)
}

func processBrowserItems(items []*analytics.AnalyticsItem, browser, browserVersion string, previousFilters []string, hasMultipleItems bool) []*AnalyticsItemWithIcon {
	result := make([]*AnalyticsItemWithIcon, len(items))
	for i, item := range items {
		icon := getBrowserIcon(item.Value)
		label := getLabel(item, browserVersion, previousFilters, false)
		isClickable := item.Drillable > 0 || hasMultipleItems

		filterKey := "b"
		filterValue := item.Value
		if browser != "" {
			filterKey = "bv"
			filterValue = getFilterValue(item, browserVersion, previousFilters)
		}

		result[i] = &AnalyticsItemWithIcon{
			AnalyticsItem: item,
			Icon:          icon,
			Label:         label,
			IsClickable:   isClickable,
			FilterKey:     filterKey,
			FilterValue:   filterValue,
		}
	}
	return result
}

func processOSItems(items []*analytics.AnalyticsItem, os, osVersion string, previousFilters []string, hasMultipleItems bool) []*AnalyticsItemWithIcon {
	result := make([]*AnalyticsItemWithIcon, len(items))
	for i, item := range items {
		label := getLabel(item, osVersion, previousFilters, false)
		isClickable := item.Drillable > 0 || hasMultipleItems

		filterKey := "os"
		filterValue := item.Value
		if os != "" {
			filterKey = "osv"
			filterValue = getFilterValue(item, osVersion, previousFilters)
		}

		result[i] = &AnalyticsItemWithIcon{
			AnalyticsItem: item,
			Icon:          "",
			Label:         label,
			IsClickable:   isClickable,
			FilterKey:     filterKey,
			FilterValue:   filterValue,
		}
	}
	return result
}

func processPageItems(items []*analytics.AnalyticsItem, page string, previousFilters []string, hasMultipleItems bool) []*AnalyticsItemWithIcon {
	result := make([]*AnalyticsItemWithIcon, len(items))
	for i, item := range items {
		label := getLabel(item, "", previousFilters, true)
		formatted := formatPageURL(item.Value)
		isClickable := item.Drillable > 0 || hasMultipleItems

		result[i] = &AnalyticsItemWithIcon{
			AnalyticsItem:  item,
			Label:          label,
			FormattedValue: formatted,
			IsClickable:    isClickable,
			FilterKey:      "pg",
			FilterValue:    item.Value,
		}
	}
	return result
}

func processReferrerItems(items []*analytics.AnalyticsItem, referrer, referrerPath string, previousFilters []string, hasMultipleItems bool) []*AnalyticsItemWithIcon {
	result := make([]*AnalyticsItemWithIcon, len(items))
	for i, item := range items {
		label := getLabel(item, referrerPath, previousFilters, true)
		faviconURL := fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=16", item.Value)
		isClickable := item.Drillable > 0 || hasMultipleItems

		filterKey := "r"
		filterValue := item.Value
		if referrer != "" {
			filterKey = "rfp"
			filterValue = getFilterValue(item, referrerPath, previousFilters)
		}

		result[i] = &AnalyticsItemWithIcon{
			AnalyticsItem: item,
			Label:         label,
			FaviconURL:    faviconURL,
			IsClickable:   isClickable,
			FilterKey:     filterKey,
			FilterValue:   filterValue,
		}
	}
	return result
}

func processCountryItems(items []*analytics.AnalyticsItem, country string, previousFilters []string, hasMultipleItems bool) []*AnalyticsItemWithIcon {
	result := make([]*AnalyticsItemWithIcon, len(items))
	for i, item := range items {
		label := getLabel(item, "", previousFilters, true)
		countryName := getCountryName(item.Value)
		isClickable := item.Drillable > 0 || hasMultipleItems

		result[i] = &AnalyticsItemWithIcon{
			AnalyticsItem: item,
			Label:         label,
			CountryName:   countryName,
			CountryCode:   strings.ToLower(item.Value),
			IsClickable:   isClickable,
			FilterKey:     "c",
			FilterValue:   item.Value,
		}
	}
	return result
}

func getBrowserIcon(browser string) string {
	lower := strings.ToLower(browser)
	if strings.Contains(lower, "retrozilla") {
		return "retrozilla"
	} else if strings.Contains(lower, "otter") {
		return "otter"
	} else if strings.Contains(lower, "dreamkey") {
		return "dreamkey"
	} else if strings.Contains(lower, "dreampassport") {
		return "dreampassport"
	} else if strings.Contains(lower, "planetweb") {
		return "planetweb"
	} else if strings.Contains(lower, "mosaic") {
		return "mosaic"
	} else if strings.Contains(lower, "netscape") {
		return "netscape"
	} else if strings.Contains(lower, "netpositive") {
		return "netpositive"
	} else if strings.Contains(lower, "konqueror") {
		return "konqueror"
	} else if strings.Contains(lower, "maxthon") {
		return "maxthon"
	} else if strings.Contains(lower, "flock") {
		return "flock"
	} else if strings.Contains(lower, "brave") {
		return "brave"
	} else if strings.Contains(lower, "chromium") {
		return "chromium"
	} else if strings.Contains(lower, "chrome") {
		return "chrome"
	} else if strings.Contains(lower, "firefox") {
		return "firefox"
	} else if strings.Contains(lower, "safari") {
		return "safari"
	} else if strings.Contains(lower, "edge") {
		return "edge"
	} else if strings.Contains(lower, "opera") {
		return "opera"
	} else if strings.Contains(lower, "ie") || strings.Contains(lower, "internet explorer") {
		return "ie"
	} else if strings.Contains(lower, "mozilla") {
		return "mozilla"
	}
	return "unknown"
}

func getLabel(item *analytics.AnalyticsItem, secondaryFilter string, previousFilters []string, showSelfWhenEmpty bool) string {
	if secondaryFilter != "" && len(previousFilters) > 1 {
		parts := append(previousFilters[1:], item.Value)
		filtered := []string{}
		for _, p := range parts {
			if p != "" {
				filtered = append(filtered, p)
			}
		}
		return strings.Join(filtered, ".")
	}

	if item.Value == "" && showSelfWhenEmpty {
		if len(previousFilters) > 0 {
			return previousFilters[len(previousFilters)-1]
		}
	}

	if item.Value == "" && !showSelfWhenEmpty {
		return "(unknown)"
	}

	return item.Value
}

func getFilterValue(item *analytics.AnalyticsItem, secondaryFilter string, previousFilters []string) string {
	if secondaryFilter != "" && len(previousFilters) >= 1 {
		parts := append(previousFilters[1:], item.Value)
		filtered := []string{}
		for _, p := range parts {
			if p != "" {
				filtered = append(filtered, p)
			}
		}
		return strings.Join(filtered, "/")
	}
	return item.Value
}

func formatPageURL(input string) string {
	if input == "" {
		return input
	}

	urlStr := input
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + input
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return input
	}

	result := parsedURL.Path
	if parsedURL.RawQuery != "" {
		result += "?" + parsedURL.RawQuery
	}

	return result
}

func getCountryName(code string) string {
	countries := map[string]string{
		"US": "United States", "GB": "United Kingdom", "CA": "Canada",
		"AU": "Australia", "DE": "Germany", "FR": "France", "IT": "Italy",
		"ES": "Spain", "NL": "Netherlands", "SE": "Sweden", "NO": "Norway",
		"DK": "Denmark", "FI": "Finland", "PL": "Poland", "RU": "Russia",
		"CN": "China", "JP": "Japan", "KR": "South Korea", "IN": "India",
		"BR": "Brazil", "MX": "Mexico", "AR": "Argentina",
	}
	if name, ok := countries[code]; ok {
		return name
	}
	return code
}

// JSONItems converts items to JSON for use in JavaScript - template function
func JSONItems(items interface{}) template.JS {
	if items == nil {
		return template.JS("[]")
	}

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		data, err := json.Marshal(items)
		if err != nil {
			return template.JS("[]")
		}
		return template.JS(data)
	}

	simpleItems := make([]map[string]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		itemMap := make(map[string]interface{})

		valueField := elem.FieldByName("Value")
		if valueField.IsValid() {
			itemMap["value"] = valueField.String()
		} else if analyticsItem := elem.FieldByName("AnalyticsItem"); analyticsItem.IsValid() {
			if analyticsItem.Kind() == reflect.Ptr {
				analyticsItem = analyticsItem.Elem()
			}
			if valueField := analyticsItem.FieldByName("Value"); valueField.IsValid() {
				itemMap["value"] = valueField.String()
			}
		}

		countField := elem.FieldByName("Count")
		if countField.IsValid() {
			itemMap["count"] = countField.Int()
		} else if analyticsItem := elem.FieldByName("AnalyticsItem"); analyticsItem.IsValid() {
			if analyticsItem.Kind() == reflect.Ptr {
				analyticsItem = analyticsItem.Elem()
			}
			if countField := analyticsItem.FieldByName("Count"); countField.IsValid() {
				itemMap["count"] = countField.Int()
			}
		}

		drillableField := elem.FieldByName("Drillable")
		if drillableField.IsValid() {
			itemMap["drillable"] = drillableField.Int()
		} else if analyticsItem := elem.FieldByName("AnalyticsItem"); analyticsItem.IsValid() {
			if analyticsItem.Kind() == reflect.Ptr {
				analyticsItem = analyticsItem.Elem()
			}
			if drillableField := analyticsItem.FieldByName("Drillable"); drillableField.IsValid() {
				itemMap["drillable"] = drillableField.Int()
			}
		}

		simpleItems[i] = itemMap
	}

	data, err := json.Marshal(simpleItems)
	if err != nil {
		return template.JS("[]")
	}
	return template.JS(data)
}
