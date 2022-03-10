package event

import (
	"fmt"
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
	databaseFileName, err := helpers.GetDatabaseFileName(item.Domain)
	if err != nil {
		panic(err)
	}

	if crawlerdetect.IsCrawler(item.UserAgent) {
		fmt.Println("crwler detected")
		return
	}

	database := db.Database{}
	database.Connect(databaseFileName)
	defer database.Close()

	// c := exec.Command("clear")
	// c.Stdout = os.Stdout
	// c.Run()
	// fmt.Println(eventQueue.GetSize())
	fmt.Println("processing", item)

	userIdent := GetSessionUserIdent(item)

	result := ua.ParseUA(item.UserAgent)

	country := geo.GetGeo(item.IP)

	session := database.GetUserSession(userIdent)

	if session == nil {
		session = database.StartUserSession(&db.UserSession{
			ID:           GetSessionId(item, item.Time),
			UserIdent:    userIdent,
			Browser:      result.Browser,
			BrowserMajor: result.BrowserMajor,
			BrowserMinor: result.BrowserMinor,
			BrowserPatch: result.BrowserPatch,
			OS:           result.OS,
			OSMajor:      result.OSMajor,
			OSMinor:      result.OSMinor,
			OSPatch:      result.OSPatch,
			Country:      country,
			SessionStart: item.Time,
			SessionEnd:   item.Time,
			UserAgent:    item.UserAgent,
			Events:       0,
			ScreenWidth:  item.ScreenWidth,
		})
	}

	session.SessionEnd = item.Time
	session.Events++

	database.UpdateUserSession(session)

	database.SaveEvent(&db.UserEvent{
		ID:        uuid.NewString(),
		Page:      item.Page,
		Name:      item.Name,
		EventTime: item.Time,
	}, session.ID)
}
