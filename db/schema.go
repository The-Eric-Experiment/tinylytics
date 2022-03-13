package db

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	UserIdent    string `gorm:"index"`
	Browser      string `gorm:"index"`
	BrowserMajor string
	BrowserMinor string
	BrowserPatch string
	OS           string `gorm:"index"`
	OSMajor      string
	OSMinor      string
	OSPatch      string
	Country      string `gorm:"index"`
	UserAgent    string
	SessionStart time.Time
	SessionEnd   time.Time
	ScreenWidth  int64
	Events       int64
}

type UserEvent struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"index"`
	Page      string
	EventTime time.Time
	SessionID string
	Session   UserSession `gorm:"foreignKey:SessionID;references:ID"`
}

type QueryFilters struct {
	Browser      *string
	BrowserMajor *string
	BrowserMinor *string
}
