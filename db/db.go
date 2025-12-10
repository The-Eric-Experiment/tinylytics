package db

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	"tinylytics/constants"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/marcboeker/go-duckdb" // DuckDB driver for database/sql
)

type Database struct {
	duckdb *sql.DB // DuckDB with raw database/sql
}

func getFilterValue(input string) string {
	if input == "null" {
		return ""
	}
	return input
}

// buildFilters builds WHERE conditions and args for raw SQL queries
func buildFilters(c *gin.Context, usePageFilter bool) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}

	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")
	os, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")
	country, hasCountry := c.GetQuery("c")
	period, hasPeriod := c.GetQuery("p")
	page, hasPage := c.GetQuery("pg")
	referer, hasReferer := c.GetQuery("r")
	refererFullPath, hasRefererFullPath := c.GetQuery("rfp")

	if !hasPeriod {
		period = constants.DATE_RAGE_24H
	}

	start, end := helpers.GetTimePeriod(period, "Australia/Sydney")

	conditions = append(conditions, "user_sessions.session_start >= ?")
	args = append(args, start)

	if end != nil {
		conditions = append(conditions, "user_sessions.session_start <= ?")
		args = append(args, end)
	}

	if hasBrowser {
		conditions = append(conditions, "user_sessions.browser = ?")
		args = append(args, getFilterValue(browser))

		if hasBrowserVersion {
			bver := strings.Split(browserVersion, "/")

			conditions = append(conditions, "user_sessions.browser_major = ?")
			args = append(args, getFilterValue(bver[0]))

			if len(bver) >= 2 {
				conditions = append(conditions, "user_sessions.browser_minor = ?")
				args = append(args, getFilterValue(bver[1]))

				if len(bver) >= 3 {
					conditions = append(conditions, "user_sessions.browser_patch = ?")
					args = append(args, getFilterValue(bver[2]))
				}
			}
		}
	}

	if hasOS {
		conditions = append(conditions, "user_sessions.os = ?")
		args = append(args, getFilterValue(os))

		if hasOSVersion {
			osver := strings.Split(osVersion, "/")

			conditions = append(conditions, "user_sessions.os_major = ?")
			args = append(args, getFilterValue(osver[0]))

			if len(osver) >= 2 {
				conditions = append(conditions, "user_sessions.os_minor = ?")
				args = append(args, getFilterValue(osver[1]))

				if len(osver) >= 3 {
					conditions = append(conditions, "user_sessions.os_patch = ?")
					args = append(args, getFilterValue(osver[2]))
				}
			}
		}
	}

	if hasCountry {
		conditions = append(conditions, "user_sessions.country = ?")
		args = append(args, getFilterValue(country))
	}

	if hasReferer {
		conditions = append(conditions, "user_sessions.referer = ?")
		args = append(args, getFilterValue(referer))

		if hasRefererFullPath {
			conditions = append(conditions, "user_sessions.referer_full_path = ?")
			args = append(args, getFilterValue(refererFullPath))
		}
	}

	if hasPage && usePageFilter {
		conditions = append(conditions, "user_events.page = ?")
		args = append(args, page)
	}

	return conditions, args
}

func (d *Database) Connect(file string) {
	// Generate DuckDB filename from database filename
	duckdbFile := file

	// Connect to DuckDB using raw database/sql
	duckdbDB, err := sql.Open("duckdb", duckdbFile)
	if err != nil {
		panic("failed to connect to DuckDB database: " + err.Error())
	}

	// Set DuckDB connection pool settings
	duckdbDB.SetMaxOpenConns(10)
	duckdbDB.SetMaxIdleConns(5)
	duckdbDB.SetConnMaxLifetime(time.Hour)

	d.duckdb = duckdbDB
}

func (d *Database) Close() {
	if d.duckdb != nil {
		d.duckdb.Close()
	}
}

func (d *Database) Initialize() {
	// Create DuckDB schema using raw SQL
	_, err := d.duckdb.Exec(`
		CREATE TABLE IF NOT EXISTS user_sessions (
			id VARCHAR PRIMARY KEY,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			user_ident VARCHAR,
			browser VARCHAR,
			browser_major VARCHAR,
			browser_minor VARCHAR,
			browser_patch VARCHAR,
			os VARCHAR,
			os_major VARCHAR,
			os_minor VARCHAR,
			os_patch VARCHAR,
			country VARCHAR,
			user_agent VARCHAR,
			referer VARCHAR,
			referer_full_path VARCHAR,
			session_start TIMESTAMP,
			session_end TIMESTAMP,
			screen_width BIGINT,
			events BIGINT
		)
	`)
	if err != nil {
		log.Printf("DuckDB user_sessions table creation failed: %v", err)
		panic("failed to create DuckDB user_sessions table")
	}

	_, err = d.duckdb.Exec(`
		CREATE TABLE IF NOT EXISTS user_events (
			id VARCHAR PRIMARY KEY,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			name VARCHAR,
			page VARCHAR,
			event_time TIMESTAMP,
			session_id VARCHAR
		)
	`)
	if err != nil {
		log.Printf("DuckDB user_events table creation failed: %v", err)
		panic("failed to create DuckDB user_events table")
	}

	// Data cleanup
	d.duckdb.Exec("update user_sessions set referer = '(none)' where referer = ''")
}

func (d *Database) GetUserSession(userIdent string) *UserSession {
	now := time.Now().UTC()
	minutes := time.Duration(-30) * time.Minute
	sessionEnd := now.Add(minutes)

	// Use raw SQL for DuckDB
	query := `
		SELECT id, created_at, updated_at, user_ident, browser, browser_major, browser_minor, 
		       browser_patch, os, os_major, os_minor, os_patch, country, user_agent, 
		       referer, referer_full_path, session_start, session_end, screen_width, events
		FROM user_sessions 
		WHERE user_ident = ? AND session_end >= ?
		LIMIT 1
	`

	var session UserSession
	err := d.duckdb.QueryRow(query, userIdent, sessionEnd).Scan(
		&session.ID, &session.CreatedAt, &session.UpdatedAt, &session.UserIdent,
		&session.Browser, &session.BrowserMajor, &session.BrowserMinor, &session.BrowserPatch,
		&session.OS, &session.OSMajor, &session.OSMinor, &session.OSPatch,
		&session.Country, &session.UserAgent, &session.Referer, &session.RefererFullPath,
		&session.SessionStart, &session.SessionEnd, &session.ScreenWidth, &session.Events,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Printf("Error fetching user session: %v", err)
		return nil
	}

	return &session
}

func (d *Database) StartUserSession(item *UserSession) *UserSession {
	// Insert into DuckDB using raw SQL
	insertSQL := `
		INSERT INTO user_sessions (
			id, created_at, updated_at, user_ident, browser, browser_major, browser_minor,
			browser_patch, os, os_major, os_minor, os_patch, country, user_agent,
			referer, referer_full_path, session_start, session_end, screen_width, events
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.duckdb.Exec(insertSQL,
		item.ID, item.CreatedAt, item.UpdatedAt, item.UserIdent, item.Browser, item.BrowserMajor, item.BrowserMinor,
		item.BrowserPatch, item.OS, item.OSMajor, item.OSMinor, item.OSPatch, item.Country, item.UserAgent,
		item.Referer, item.RefererFullPath, item.SessionStart, item.SessionEnd, item.ScreenWidth, item.Events,
	)

	if err != nil {
		log.Printf("ERROR: Failed to write session to DuckDB: %v", err)
		panic(err)
	}

	return item
}

func (d *Database) UpdateUserSession(item *UserSession) {
	// Update in DuckDB using raw SQL
	updateSQL := `
		UPDATE user_sessions SET
			created_at = ?, updated_at = ?, user_ident = ?, browser = ?, browser_major = ?,
			browser_minor = ?, browser_patch = ?, os = ?, os_major = ?, os_minor = ?,
			os_patch = ?, country = ?, user_agent = ?, referer = ?, referer_full_path = ?,
			session_start = ?, session_end = ?, screen_width = ?, events = ?
		WHERE id = ?
	`

	_, err := d.duckdb.Exec(updateSQL,
		item.CreatedAt, item.UpdatedAt, item.UserIdent, item.Browser, item.BrowserMajor,
		item.BrowserMinor, item.BrowserPatch, item.OS, item.OSMajor, item.OSMinor,
		item.OSPatch, item.Country, item.UserAgent, item.Referer, item.RefererFullPath,
		item.SessionStart, item.SessionEnd, item.ScreenWidth, item.Events,
		item.ID,
	)

	if err != nil {
		log.Printf("ERROR: Failed to update session in DuckDB: %v", err)
		panic(err)
	}
}

func (d *Database) SaveEvent(item *UserEvent, sessionId string) *UserEvent {
	// Insert into DuckDB using raw SQL
	insertSQL := `
		INSERT INTO user_events (
			id, created_at, updated_at, name, page, event_time, session_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.duckdb.Exec(insertSQL,
		item.ID, item.CreatedAt, item.UpdatedAt, item.Name, item.Page, item.EventTime, sessionId,
	)

	if err != nil {
		log.Printf("ERROR: Failed to write event to DuckDB: %v", err)
		panic(err)
	}

	return item
}

func (d *Database) GetSessions(c *gin.Context) int64 {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM user_sessions 
		WHERE %s
	`, strings.Join(conditions, " AND "))

	var count int64
	err := d.duckdb.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Printf("ERROR: Failed to get sessions count: %v", err)
		return 0
	}

	return count
}

func (d *Database) GetPageViews(c *gin.Context) int64 {
	conditions, args := buildFilters(c, false)

	// Add pageview condition
	allConditions := append([]string{"user_events.name = ?"}, conditions...)
	allArgs := append([]interface{}{"pageview"}, args...)

	// Add page filter if present
	page, hasPage := c.GetQuery("pg")
	if hasPage {
		allConditions = append(allConditions, "user_events.page = ?")
		allArgs = append(allArgs, page)
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM user_events 
		LEFT JOIN user_sessions ON user_sessions.id = user_events.session_id 
		WHERE %s
	`, strings.Join(allConditions, " AND "))

	var count int64
	err := d.duckdb.QueryRow(query, allArgs...).Scan(&count)
	if err != nil {
		log.Printf("ERROR: Failed to get page views count: %v", err)
		return 0
	}

	return count
}

func (d *Database) GetAvgSessionDuration(c *gin.Context) float64 {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	query := fmt.Sprintf(`
		SELECT AVG(EXTRACT(EPOCH FROM (session_end - session_start)))
		FROM user_sessions 
		WHERE %s
	`, strings.Join(conditions, " AND "))

	var duration sql.NullFloat64
	err := d.duckdb.QueryRow(query, args...).Scan(&duration)
	if err != nil {
		log.Printf("ERROR: Failed to get avg session duration: %v", err)
		return 0
	}

	if !duration.Valid {
		return 0
	}

	return duration.Float64
}

func (d *Database) GetBounceRate(c *gin.Context) int64 {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	query := fmt.Sprintf(`
		SELECT 
			SUM(CASE WHEN EXTRACT(EPOCH FROM (session_end - session_start)) = 0.0 THEN 1 ELSE 0 END) as bounces,
			COUNT(*) as total
		FROM user_sessions 
		WHERE %s
	`, strings.Join(conditions, " AND "))

	var bounces, total sql.NullFloat64
	err := d.duckdb.QueryRow(query, args...).Scan(&bounces, &total)
	if err != nil {
		log.Printf("ERROR: Failed to get bounce rate: %v", err)
		return 0
	}

	if !bounces.Valid || !total.Valid || total.Float64 == 0 {
		return 0
	}

	return int64(math.Round((bounces.Float64 / total.Float64) * 100))
}

func (d *Database) GetBrowsers(c *gin.Context) (*sql.Rows, error) {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	_, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	var query string

	if !hasBrowser {
		// Base browser query
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.browser as value,
				COUNT(user_sessions.browser) as count,
				SUM(CASE WHEN user_sessions.browser_major <> '' AND user_sessions.browser_major <> '0' THEN 1 ELSE 0 END) AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.browser
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	} else if !hasBrowserVersion {
		// Browser major version
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.browser_major as value,
				COUNT(user_sessions.browser_major) as count,
				SUM(CASE WHEN user_sessions.browser_minor <> '' AND user_sessions.browser_minor <> '0' THEN 1 ELSE 0 END) AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.browser_major
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	} else {
		bver := strings.Split(browserVersion, "/")
		if len(bver) < 2 {
			// Browser minor version
			query = fmt.Sprintf(`
				SELECT 
					user_sessions.browser_minor as value,
					COUNT(user_sessions.browser_minor) as count,
					SUM(CASE WHEN user_sessions.browser_patch <> '' AND user_sessions.browser_patch <> '0' THEN 1 ELSE 0 END) AS drillable
				FROM user_sessions 
				WHERE %s
				GROUP BY user_sessions.browser_minor
				ORDER BY count DESC
				LIMIT 20
			`, strings.Join(conditions, " AND "))
		} else {
			// Browser patch version
			query = fmt.Sprintf(`
				SELECT 
					user_sessions.browser_patch as value,
					COUNT(user_sessions.browser_patch) as count,
					0 AS drillable
				FROM user_sessions 
				WHERE %s
				GROUP BY user_sessions.browser_patch
				ORDER BY count DESC
				LIMIT 20
			`, strings.Join(conditions, " AND "))
		}
	}

	return d.duckdb.Query(query, args...)
}

func (d *Database) GetOSs(c *gin.Context) (*sql.Rows, error) {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	_, hasOS := c.GetQuery("os")
	osVersion, hasOSVersion := c.GetQuery("osv")

	var query string

	if !hasOS {
		// Base OS query
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.os as value,
				COUNT(user_sessions.os) as count,
				SUM(CASE WHEN user_sessions.os_major <> '' AND user_sessions.os_major <> '0' THEN 1 ELSE 0 END) AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.os
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	} else if !hasOSVersion {
		// OS major version
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.os_major as value,
				COUNT(user_sessions.os_major) as count,
				SUM(CASE WHEN user_sessions.os_minor <> '' AND user_sessions.os_minor <> '0' THEN 1 ELSE 0 END) AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.os_major
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	} else {
		osver := strings.Split(osVersion, "/")
		if len(osver) < 2 {
			// OS minor version
			query = fmt.Sprintf(`
				SELECT 
					user_sessions.os_minor as value,
					COUNT(user_sessions.os_minor) as count,
					SUM(CASE WHEN user_sessions.os_patch <> '' AND user_sessions.os_patch <> '0' THEN 1 ELSE 0 END) AS drillable
				FROM user_sessions 
				WHERE %s
				GROUP BY user_sessions.os_minor
				ORDER BY count DESC
				LIMIT 20
			`, strings.Join(conditions, " AND "))
		} else {
			// OS patch version
			query = fmt.Sprintf(`
				SELECT 
					user_sessions.os_patch as value,
					COUNT(user_sessions.os_patch) as count,
					0 AS drillable
				FROM user_sessions 
				WHERE %s
				GROUP BY user_sessions.os_patch
				ORDER BY count DESC
				LIMIT 20
			`, strings.Join(conditions, " AND "))
		}
	}

	return d.duckdb.Query(query, args...)
}

func (d *Database) GetCountries(c *gin.Context) (*sql.Rows, error) {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	query := fmt.Sprintf(`
		SELECT 
			user_sessions.country as value,
			COUNT(user_sessions.country) as count,
			0 AS drillable
		FROM user_sessions 
		WHERE %s
		GROUP BY user_sessions.country
		ORDER BY count DESC
	`, strings.Join(conditions, " AND "))

	return d.duckdb.Query(query, args...)
}

func (d *Database) GetReferrers(c *gin.Context) (*sql.Rows, error) {
	conditions, args := buildFilters(c, false) // No user_events table, so no page filter

	_, hasReferrer := c.GetQuery("r")

	var query string

	if !hasReferrer {
		// Base referrer query
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.referer as value,
				COUNT(user_sessions.referer) as count,
				SUM(CASE WHEN user_sessions.referer_full_path <> '' THEN 1 ELSE 0 END) AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.referer
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	} else {
		// Referrer full path
		query = fmt.Sprintf(`
			SELECT 
				user_sessions.referer_full_path as value,
				COUNT(user_sessions.referer_full_path) as count,
				0 AS drillable
			FROM user_sessions 
			WHERE %s
			GROUP BY user_sessions.referer_full_path
			ORDER BY count DESC
			LIMIT 20
		`, strings.Join(conditions, " AND "))
	}

	return d.duckdb.Query(query, args...)
}

func (d *Database) GetPages(c *gin.Context) (*sql.Rows, error) {
	conditions, args := buildFilters(c, false)

	// Add pageview condition
	allConditions := append([]string{"user_events.name = ?"}, conditions...)
	allArgs := append([]interface{}{"pageview"}, args...)

	// Add page filter if present
	page, hasPage := c.GetQuery("pg")
	if hasPage {
		allConditions = append(allConditions, "user_events.page = ?")
		allArgs = append(allArgs, page)
	}

	query := fmt.Sprintf(`
		SELECT 
			user_events.page as value,
			COUNT(user_events.page) as count,
			0 AS drillable
		FROM user_events 
		LEFT JOIN user_sessions ON user_sessions.id = user_events.session_id 
		WHERE %s
		GROUP BY user_events.page
		ORDER BY count DESC
		LIMIT 20
	`, strings.Join(allConditions, " AND "))

	return d.duckdb.Query(query, allArgs...)
}
