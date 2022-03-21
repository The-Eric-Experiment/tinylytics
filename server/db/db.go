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
		Logger: logger.Default.LogMode(logger.Silent),
		// Logger: logger.Default.LogMode(logger.Info),
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

	d.db.Exec("update user_sessions set referer = '(none)' where referer = ''")
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

func setFilters(db *gorm.DB, c *gin.Context) *gorm.DB {
	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")
	os, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")
	country, hasCountry := c.GetQuery("c")
	period, hasPeriod := c.GetQuery("p")
	referer, hasReferer := c.GetQuery("r")

	if !hasPeriod {
		period = constants.DATE_RAGE_24H
	}

	start, end := helpers.GetTimePeriod(period, "Australia/Sydney")

	db = db.Where("user_sessions.session_start >= ?", start)

	if end != nil {
		db = db.Where("user_sessions.session_start <= ?", end)
	}

	if hasBrowser {
		db = db.Where(&UserSession{Browser: browser})

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")
			bmj := bver[0]

			db = db.Where(&UserSession{BrowserMajor: bmj})
			if len(bver) >= 2 {
				db = db.Where(&UserSession{BrowserMinor: bver[1]})
			}
		}
	}

	if hasOS {
		db = db.Where(&UserSession{OS: os})

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")
			osmj := osver[0]

			db = db.Where(&UserSession{OSMajor: osmj})
			if len(osver) >= 2 {
				db = db.Where(&UserSession{OSMinor: osver[1]})
			}
		}
	}

	if hasCountry {
		db = db.Where(&UserSession{Country: country})
	}

	if hasReferer {
		db = db.Where(&UserSession{Referer: referer})
	}

	return db
}

func (d *Database) GetSessions(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserSession{})

	q = setFilters(q, c)

	q.Count(&count)
	return count
}

func (d *Database) GetPageViews(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserEvent{}).Joins("left join user_sessions on user_sessions.id = user_events.session_id").Where(&UserEvent{Name: "pageview"})

	q = setFilters(q, c)

	q.Count(&count)
	return count
}

func (d *Database) GetBrowsers(c *gin.Context) (*sql.Rows, error) {
	querySelect := "browser as name, count(browser) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Group("browser").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = setFilters(q, c).Limit(20)

	_, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	if hasBrowser {
		q = q.Group("user_sessions.browser_major")

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")

			q = q.Group("user_sessions.browser_minor")
			if len(bver) >= 2 {
				q = q.Group("user_sessions.browser_patch")
			}
		}
	}

	if hasBrowser {
		querySelect += ", user_sessions.browser_major as major"
		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")
			querySelect += ", user_sessions.browser_minor as minor"
			if len(bver) >= 2 {
				querySelect += ", user_sessions.browser_patch as patch"
			}
		}
	}

	q.Select(querySelect)

	return q.Rows()
}

func (d *Database) GetOSs(c *gin.Context) (*sql.Rows, error) {
	querySelect := "user_sessions.os as name, count(user_sessions.os) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Group("os").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = setFilters(q, c).Limit(20)
	_, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")

	if hasOS {
		q = q.Group("user_sessions.os_major")

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")

			q = q.Group("user_sessions.os_minor")
			if len(osver) >= 2 {
				q = q.Group("user_sessions.os_patch")
			}
		}
	}

	if hasOS {
		querySelect += ", user_sessions.os_major as major"

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")
			querySelect += ", user_sessions.os_minor as minor"
			if len(osver) >= 2 {
				querySelect += ", user_sessions.os_patch as patch"
			}
		}
	}

	q.Select(querySelect)

	return q.Rows()
}

func (d *Database) GetCountries(c *gin.Context) (*sql.Rows, error) {
	querySelect := "country as name, count(country) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	}).Group("country").Limit(20)

	q = setFilters(q, c)

	return q.Rows()
}

func (d *Database) GetReferrers(c *gin.Context) (*sql.Rows, error) {
	querySelect := "referer as name, count(referer) as count"

	q := d.db.Model(&UserSession{}).Select(querySelect).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	}).Group("referer").Limit(20)

	q = setFilters(q, c)

	_, hasReferrer := c.GetQuery("r")

	if hasReferrer {
		q = q.Group("user_sessions.referer_full_path")
	}

	if hasReferrer {
		querySelect += ", user_sessions.referer_full_path as major"
	}

	q.Select(querySelect)

	return q.Rows()
}

func (d *Database) Scan(rows *sql.Rows, dest interface{}) error {
	return d.db.ScanRows(rows, &dest)
}
