package db

import (
	"database/sql"
	"errors"
	"time"

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

func (d *Database) GetBrowsersQuery(filters *QueryFilters) (*sql.Rows, error) {
	querySelect := "browser as name, count(browser) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Group("browser").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	if filters.Browser != nil {
		q = q.Where(&UserSession{Browser: *filters.Browser}).Group("browser_major")

		querySelect += ", browser_major"
		if filters.BrowserMajor != nil {
			q = q.Where(&UserSession{BrowserMajor: *filters.BrowserMajor}).Group("browser_minor")
			querySelect += ", browser_minor"
		}
	}

	q = q.Select(querySelect)

	return q.Rows()
}
