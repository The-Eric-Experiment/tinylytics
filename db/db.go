package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"tinylytics/constants"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func (d *Database) Connect(file string) {
	db, err := gorm.Open(sqlite.Open("file:"+file+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	d.db = db
}

func (d *Database) Close() {
	db, err := d.db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	db.Close()
}

func (d *Database) Initialize() {
	// Migrate the schema
	d.db.AutoMigrate(&UserSession{})
	d.db.AutoMigrate(&UserEvent{})
}

func (d *Database) GetUserSession(userIdent string) *UserSession {
	now := time.Now().UTC()
	minutes := time.Duration(-30) * time.Minute
	sessionEnd := now.Add(minutes)

	var session *UserSession = nil
	if result := d.db.Where("user_ident = ? and session_end >= ?", userIdent, sessionEnd).First(&session); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(result.Error)
	}

	return session
}

func (d *Database) StartUserSession(item *UserSession) *UserSession {
	d.db.Create(&item)
	return item
}

func (d *Database) UpdateUserSession(item *UserSession) {
	d.db.Save(&item)
}

func (d *Database) SaveEvent(item *UserEvent, sessionId string) *UserEvent {
	item.Session = UserSession{ID: sessionId}
	d.db.Create(&item)
	return item
}

func (d *Database) GetSessions(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserSession{})

	q = setFilters(q, c, "")

	q.Count(&count)
	return count
}

func (d *Database) GetPageViews(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserEvent{}).Joins("left join user_sessions on user_sessions.id = user_events.session_id").Where(&UserEvent{Name: "pageview"})

	q = setFilters(q, c, "")

	q.Count(&count)
	return count
}

func setFilters(db *gorm.DB, c *gin.Context, querySelect string) *gorm.DB {
	var qs string = ""
	if querySelect != "" {
		qs = querySelect
	}
	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")
	period, hasPeriod := c.GetQuery("p")

	if !hasPeriod {
		period = constants.DATE_RAGE_24H
	}

	start, end := helpers.GetTimePeriod(period, "Australia/Sydney")

	db = db.Where("user_sessions.session_start >= ?", start)

	if end != nil {
		db = db.Where("user_sessions.session_start <= ?", end)
	}

	if hasBrowser {
		db = db.Where(&UserSession{Browser: browser}).Group("user_sessions.browser_major")
		qs += ", user_sessions.browser_major"

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, ".")
			bmj := bver[0]

			db = db.Where(&UserSession{BrowserMajor: bmj}).Group("user_sessions.browser_minor")
			qs += ", user_sessions.browser_minor"
			if len(bver) >= 2 {
				db = db.Where(&UserSession{BrowserMinor: bver[1]}).Group("user_sessions.browser_patch")
				qs += ", user_sessions.browser_patch"
			}

		}

	}

	if querySelect != "" {
		db.Select(querySelect)
	}

	return db
}

func (d *Database) GetBrowsers(c *gin.Context) (*sql.Rows, error) {
	querySelect := "browser as name, count(browser) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Group("browser").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = setFilters(q, c, querySelect).Limit(20)

	return q.Rows()
}

func (d *Database) Scan(rows *sql.Rows, dest interface{}) error {
	return d.db.ScanRows(rows, &dest)
}
