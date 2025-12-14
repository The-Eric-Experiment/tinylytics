package event

import (
	"log"
	"net/url"
	"strings"
	"time"
	"tinylytics/db"
	"tinylytics/geo"
	"tinylytics/helpers"
	"tinylytics/ua"

	"github.com/google/uuid"
	"github.com/x-way/crawlerdetect"
)

type ClientInfo struct {
	Name                      string
	UserAgent                 string
	HostName                  string
	Domain                    string
	Page                      string
	ClientHintUA              string
	ClientHintMobile          string
	ClientHintPlatform        string
	ClientHintFullVersion     string
	ClientHintPlatformVersion string
	IP                        string
	Referer                   string
	Time                      time.Time
	ScreenWidth               int64
}

type EventData struct {
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Page        string `json:"page"`
	ScreenWidth int64  `json:"screenWidth"`
}

func ProcessEvent(item *ClientInfo) {
	if crawlerdetect.IsCrawler(item.UserAgent) {
		log.Printf("[QUEUE] Crawler detected - skipping: %s", item.UserAgent)
		return
	}

	database, err := db.GetDatabaseByDomain(item.Domain)
	if err != nil {
		log.Printf("ERROR: Failed to get database for domain %s: %v", item.Domain, err)
		return
	}

	log.Printf("[QUEUE] Processing event: domain=%s page=%s IP=%s", item.Domain, item.Page, item.IP)

	userIdent := GetSessionUserIdent(item)

	result := ua.ParseUA(item.UserAgent)

	country := geo.GetGeo(item.IP)

	// Use item.Time (event timestamp) instead of processing time to correctly match sessions
	session := database.GetUserSessionAtTime(userIdent, item.Time)

	if session == nil {
		referrerDomain, referrerFullPath := helpers.FilterReferrer(item.Referer, item.Domain)

		session = database.StartUserSession(&db.UserSession{
			ID:              GetSessionId(item, item.Time),
			UserIdent:       userIdent,
			Browser:         result.Browser,
			BrowserMajor:    result.BrowserMajor,
			BrowserMinor:    result.BrowserMinor,
			BrowserPatch:    result.BrowserPatch,
			OS:              result.OS,
			OSMajor:         result.OSMajor,
			OSMinor:         result.OSMinor,
			OSPatch:         result.OSPatch,
			Country:         country,
			SessionStart:    item.Time,
			SessionEnd:      item.Time,
			UserAgent:       item.UserAgent,
			Referer:         referrerDomain,
			RefererFullPath: referrerFullPath,
			Events:          0,
			ScreenWidth:     item.ScreenWidth,
		})
	}

	session.SessionEnd = item.Time
	session.Events++

	database.UpdateUserSession(session)

	var page = item.Page

	path, err := url.Parse(page)

	if err == nil {
		page = strings.Join([]string{
			strings.Trim(item.Domain, "/"),
			strings.Trim(path.Path, "/"),
		}, "/")
	}

	database.SaveEvent(&db.UserEvent{
		ID:        uuid.NewString(),
		Page:      page,
		Name:      item.Name,
		EventTime: item.Time,
	}, session.ID)
}
