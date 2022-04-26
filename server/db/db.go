package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
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

func getFilterValue(input string) string {
	if input == "null" {
		return ""
	}
	return input
}

func setFilters(db *gorm.DB, c *gin.Context, usePageFilter bool) *gorm.DB {
	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")
	os, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")
	country, hasCountry := c.GetQuery("c")
	period, hasPeriod := c.GetQuery("p")
	page, hasPage := c.GetQuery("pg")
	referer, hasReferer := c.GetQuery("r")
	refererFullPath, hasRefererFullPath := c.GetQuery("rfp")

	fmt.Println(hasRefererFullPath)
	fmt.Println(refererFullPath)

	if !hasPeriod {
		period = constants.DATE_RAGE_24H
	}

	start, end := helpers.GetTimePeriod(period, "Australia/Sydney")

	db = db.Where("user_sessions.session_start >= ?", start)

	if end != nil {
		db = db.Where("user_sessions.session_start <= ?", end)
	}

	if hasBrowser {
		db = db.Where("user_sessions.browser = ?", getFilterValue(browser))

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")

			db = db.Where("user_sessions.browser_major = ?", getFilterValue(bver[0]))
			if len(bver) >= 2 {
				db = db.Where("user_sessions.browser_minor = ?", getFilterValue(bver[1]))
				if len(bver) >= 3 {
					db = db.Where("user_sessions.browser_patch = ?", getFilterValue(bver[2]))
				}
			}
		}
	}

	if hasOS {
		db = db.Where("user_sessions.os = ? ", getFilterValue(os))

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")

			db = db.Where("user_sessions.os_major = ?", getFilterValue(osver[0]))
			if len(osver) >= 2 {
				db = db.Where("user_sessions.os_minor = ?", getFilterValue(osver[1]))
				if len(osver) >= 3 {
					db = db.Where("user_sessions.os_patch = ?", getFilterValue(osver[2]))
				}
			}
		}
	}

	if hasCountry {
		db = db.Where("user_sessions.country = ?", getFilterValue(country))
	}

	if hasReferer {
		db = db.Where("user_sessions.referer = ?", getFilterValue(referer))
		if hasRefererFullPath {
			db = db.Where("referer_full_path = ?", getFilterValue(refererFullPath))
		}
	}

	if hasPage && usePageFilter {
		db = db.Joins("left join user_events on user_sessions.id = user_events.session_id")
		db = db.Where("user_events.page", page)
	}

	return db
}

func selector(db *gorm.DB, fields ...string) *gorm.DB {
	sel := strings.Join(fields, ", ")
	db.Select(sel)
	return db
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

func (d *Database) GetSessions(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserSession{})

	q = setFilters(q, c, true)

	q.Count(&count)
	return count
}

func (d *Database) GetPageViews(c *gin.Context) int64 {
	var count int64
	q := d.db.Model(&UserEvent{}).Joins("left join user_sessions on user_sessions.id = user_events.session_id").Where(&UserEvent{Name: "pageview"})

	q = setFilters(q, c, false)

	page, hasPage := c.GetQuery("pg")

	if hasPage {
		q = q.Where("user_events.page = ?", page)
	}

	q.Count(&count)
	return count
}

type Blah struct {
	Browser  string
	Duration int64
}

func (d *Database) GetAvgSessionDuration(c *gin.Context) float64 {
	q := d.db.Model(&UserSession{})
	q = setFilters(q, c, true)
	var duration float64
	row := q.Select(`avg((julianday(user_sessions.session_end) - julianday(user_sessions.session_start)) * 86400.0)`).Row()
	row.Scan(&duration)
	return duration
}

func (d *Database) GetBounceRate(c *gin.Context) int64 {
	q := d.db.Model(&UserSession{})
	q = setFilters(q, c, true)
	var bounces float64
	var total float64
	row := q.Select(`sum(CASE WHEN ((julianday(user_sessions.session_end) - julianday(user_sessions.session_start)) * 86400.0) = 0.0 THEN 1 ELSE 0 END) as bounces, count(*) as total`).Row()
	row.Scan(&bounces, &total)
	return int64(math.Round((bounces / total) * 100))
}

func (d *Database) GetBrowsers(c *gin.Context) (*sql.Rows, error) {
	q := d.db.Model(&UserSession{}).Group("browser").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = selector(q,
		"user_sessions.browser as value",
		"count(user_sessions.browser) as count",
		"SUM(CASE WHEN user_sessions.browser_major <> '' AND user_sessions.browser_major <> '0' THEN 1 ELSE 0 END) AS drillable",
	)

	q = setFilters(q, c, true).Limit(20)

	_, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	if hasBrowser {
		q = q.Group("user_sessions.browser_major")
		q = selector(q,
			"user_sessions.browser_major as value",
			"count(user_sessions.browser_major) as count",
			"SUM(CASE WHEN user_sessions.browser_minor <> '' AND user_sessions.browser_minor <> '0' THEN 1 ELSE 0 END) AS drillable",
		)

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")
			q = selector(q,
				"user_sessions.browser_minor as value",
				"count(user_sessions.browser_minor) as count",
				"SUM(CASE WHEN user_sessions.browser_patch <> '' AND user_sessions.browser_patch <> '0' THEN 1 ELSE 0 END) AS drillable",
			)

			q = q.Group("user_sessions.browser_minor")
			if len(bver) >= 2 {
				q = q.Group("user_sessions.browser_patch")
				q = selector(q,
					"user_sessions.browser_patch as value",
					"count(user_sessions.browser_patch) as count",
					"0 AS drillable",
				)
			}
		}
	}

	return q.Rows()
}

func (d *Database) GetOSs(c *gin.Context) (*sql.Rows, error) {

	q := d.db.Model(&UserSession{}).Group("os").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	})

	q = setFilters(q, c, true).Limit(20)

	q = selector(q,
		"user_sessions.os as value",
		"count(user_sessions.os) as count",
		"SUM(CASE WHEN user_sessions.os_major <> '' AND user_sessions.os_major <> '0' THEN 1 ELSE 0 END) AS drillable",
	)

	_, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")

	if hasOS {
		q = q.Group("user_sessions.os_major")
		q = selector(q,
			"user_sessions.os_major as value",
			"count(user_sessions.os_major) as count",
			"SUM(CASE WHEN user_sessions.os_minor <> '' AND user_sessions.os_minor <> '0' THEN 1 ELSE 0 END) AS drillable",
		)

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")
			q = q.Group("user_sessions.os_minor")
			q = selector(q,
				"user_sessions.os_minor as value",
				"count(user_sessions.os_minor) as count",
				"SUM(CASE WHEN user_sessions.os_patch <> '' AND user_sessions.os_patch <> '0' THEN 1 ELSE 0 END) AS drillable",
			)

			if len(osver) >= 2 {
				q = q.Group("user_sessions.os_patch")
				q = selector(q,
					"user_sessions.os_patch as value",
					"count(user_sessions.os_patch) as count",
					"0 AS drillable",
				)
			}
		}
	}

	return q.Rows()
}

func (d *Database) GetCountries(c *gin.Context) (*sql.Rows, error) {
	q := d.db.Model(&UserSession{}).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	}).Group("country").Limit(20)

	q = setFilters(q, c, true)

	q = selector(q,
		"user_sessions.country as value",
		"count(user_sessions.country) as count",
		"0 AS drillable",
	)

	return q.Rows()
}

func (d *Database) GetReferrers(c *gin.Context) (*sql.Rows, error) {
	q := d.db.Model(&UserSession{}).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	}).Group("referer").Limit(20)

	q = selector(q,
		"user_sessions.referer as value",
		"count(user_sessions.referer) as count",
		"SUM(CASE WHEN user_sessions.referer_full_path <> '' THEN 1 ELSE 0 END) AS drillable",
	)

	q = setFilters(q, c, true)

	_, hasReferrer := c.GetQuery("r")

	if hasReferrer {
		q = q.Group("user_sessions.referer_full_path")
	}

	if hasReferrer {
		q = selector(q,
			"user_sessions.referer_full_path as value",
			"count(user_sessions.referer_full_path) as count",
			"0 AS drillable",
		)
	}

	return q.Rows()
}

func (d *Database) GetPages(c *gin.Context) (*sql.Rows, error) {
	q := d.db.Model(&UserEvent{}).Joins("left join user_sessions on user_sessions.id = user_events.session_id").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "count desc", WithoutParentheses: true},
	}).Where(&UserEvent{Name: "pageview"}).Group("page").Limit(20)

	q = selector(q,
		"user_events.page as value",
		"count(user_events.page) as count",
		"0 AS drillable",
	)

	q = setFilters(q, c, false)

	page, hasPage := c.GetQuery("pg")

	if hasPage {
		q = q.Where("user_events.page = ?", page)
	}

	return q.Rows()
}

func (d *Database) Scan(rows *sql.Rows, dest interface{}) error {
	return d.db.ScanRows(rows, &dest)
}
