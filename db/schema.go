package db

import (
	"time"
)

// =============================================================================
// DuckDB Schema - Columnar storage optimized for analytics
// =============================================================================

type UserSession struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserIdent       string
	Browser         string
	BrowserMajor    string
	BrowserMinor    string
	BrowserPatch    string
	OS              string
	OSMajor         string
	OSMinor         string
	OSPatch         string
	Country         string
	UserAgent       string
	Referer         string
	RefererFullPath string
	SessionStart    time.Time
	SessionEnd      time.Time
	ScreenWidth     int64
	Events          int64
}

type UserEvent struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Page      string
	EventTime time.Time
	SessionID string
}

type QueryFilters struct {
	Browser      *string
	BrowserMajor *string
	BrowserMinor *string
}
