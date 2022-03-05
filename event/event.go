package event

import (
	"time"
	"tinylytics/db"
	"tinylytics/geo"
	"tinylytics/helpers"
	"tinylytics/ua"

	"github.com/google/uuid"
)

type ClientInfo struct {
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
}

type EventData struct {
	Domain string `json:"domain"`
	Page   string `json:"page"`
}

func ProcessEvent(item *ClientInfo) {
	databaseFileName := helpers.GetDatabaseFileName(item.Domain)
	database := db.Database{}
	database.Connect(databaseFileName)
	defer database.Close()

	// c := exec.Command("clear")
	// 	c.Stdout = os.Stdout
	// 	c.Run()
	// fmt.Println(eventQueue.GetSize())

	userIdent := GetSessionUserIdent(item)

	result := ua.ParseUA(item.UserAgent)

	country := geo.GetGeo(item.IP)

	session := database.GetUserSession(userIdent)

	if session == nil {
		session = database.StartUserSession(&db.UserSession{
			ID:             GetSessionId(item, item.Time),
			UserIdent:      userIdent,
			Browser:        result.Browser,
			BrowserVersion: result.BrowserVersion,
			OS:             result.OS,
			OSVersion:      result.OSVersion,
			Country:        country,
			SessionStart:   item.Time,
			SessionEnd:     item.Time,
			UserAgent:      item.UserAgent,
			Events:         0,
		})
	}

	session.SessionEnd = item.Time
	session.Events++

	database.UpdateUserSession(session)

	database.SaveEvent(&db.UserEvent{
		ID:        uuid.NewString(),
		Page:      item.Page,
		EventTime: item.Time,
	}, session.ID)
}
