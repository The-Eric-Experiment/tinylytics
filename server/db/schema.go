package db

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	ID              string `gorm:"primaryKey;index:idx_sessions_id_start,priority:1"`
	UserIdent       string `gorm:"index:idx_user_ident_session_end,priority:1"`
	Browser         string `gorm:"index:idx_sessions_start_browser,priority:2;index:idx_sessions_browser_major,priority:1"`
	BrowserMajor    string `gorm:"index:idx_sessions_browser_major,priority:2;index:idx_sessions_browser_minor,priority:1"`
	BrowserMinor    string `gorm:"index:idx_sessions_browser_minor,priority:2;index:idx_sessions_browser_patch,priority:1"`
	BrowserPatch    string `gorm:"index:idx_sessions_browser_patch,priority:2"`
	OS              string `gorm:"index:idx_sessions_start_os,priority:2;index:idx_sessions_os_major,priority:1"`
	OSMajor         string `gorm:"index:idx_sessions_os_major,priority:2;index:idx_sessions_os_minor,priority:1"`
	OSMinor         string `gorm:"index:idx_sessions_os_minor,priority:2;index:idx_sessions_os_patch,priority:1"`
	OSPatch         string `gorm:"index:idx_sessions_os_patch,priority:2"`
	Country         string `gorm:"index:idx_sessions_start_country,priority:2"`
	UserAgent       string
	Referer         string    `gorm:"index:idx_sessions_start_referer,priority:2;index:idx_sessions_referer_path,priority:1"`
	RefererFullPath string    `gorm:"index:idx_sessions_referer_path,priority:2"`
	SessionStart    time.Time `gorm:"index;index:idx_sessions_start_browser,priority:1;index:idx_sessions_start_country,priority:1;index:idx_sessions_start_os,priority:1;index:idx_sessions_start_referer,priority:1;index:idx_sessions_start_end,priority:1;index:idx_sessions_id_start,priority:2"`
	SessionEnd      time.Time `gorm:"index:idx_user_ident_session_end,priority:2;index:idx_sessions_start_end,priority:2"`
	ScreenWidth     int64
	Events          int64
}

type UserEvent struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"index:idx_events_session_name,priority:2;index:idx_events_name_session_page,priority:1"`
	Page      string `gorm:"index:idx_events_session_page,priority:2;index:idx_events_name_session_page,priority:3;index:idx_events_pageview_page"`
	EventTime time.Time
	SessionID string      `gorm:"index:idx_events_session_name,priority:1;index:idx_events_session_page,priority:1;index:idx_events_name_session_page,priority:2"`
	Session   UserSession `gorm:"foreignKey:SessionID;references:ID"`
}

type QueryFilters struct {
	Browser      *string
	BrowserMajor *string
	BrowserMinor *string
}
