package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"

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

func (d *Database) GetSessions() int64 {
	var count int64
	d.db.Model(&UserSession{}).Count(&count)

	return count
}

func (d *Database) GetPageViews() int64 {
	var count int64
	d.db.Model(&UserEvent{}).Where(&UserEvent{Name: "pageview"}).Count(&count)

	return count
}

func setBrowserFilters(db *gorm.DB, c *gin.Context, querySelect string) *gorm.DB {
	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	if hasBrowser {
		db = db.Where(&UserSession{Browser: browser}).Group("browser_major")
		querySelect += ", browser_major"

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, ".")
			bmj := bver[0]

			db = db.Where(&UserSession{BrowserMajor: bmj}).Group("browser_minor")
			querySelect += ", browser_minor"
			if len(bver) >= 2 {
				db = db.Where(&UserSession{BrowserMinor: bver[1]}).Group("browser_patch")
				querySelect += ", browser_patch"
			}

		}

	}

	db.Select(querySelect)

	return db
}

func (d *Database) GetBrowsers(c *gin.Context) (*sql.Rows, error) {
	querySelect := "browser as name, count(browser) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Group("browser").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = setBrowserFilters(q, c, querySelect)

	return q.Rows()
}

func (d *Database) Scan(rows *sql.Rows, dest interface{}) error {
	return d.db.ScanRows(rows, &dest)
}
