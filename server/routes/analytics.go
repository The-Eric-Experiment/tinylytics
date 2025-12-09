package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"tinylytics/analytics"
	"tinylytics/config"
	"tinylytics/db"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
)

func arrayFromRows(rows *sql.Rows) []*analytics.AnalyticsItem {
	if rows == nil {
		log.Printf("ERROR: rows is nil in arrayFromRows")
		return make([]*analytics.AnalyticsItem, 0)
	}
	defer rows.Close()

	list := make([]*analytics.AnalyticsItem, 0)

	for rows.Next() {
		var output analytics.AnalyticsItem

		// All analytics queries return: value, count, drillable
		err := rows.Scan(&output.Value, &output.Count, &output.Drillable)
		if err != nil {
			log.Printf("ERROR: Failed to scan row: %v", err)
			continue
		}

		list = append(list, &output)
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Error iterating rows: %v", err)
	}

	return list
}

func getDB(c *gin.Context) *db.Database {
	domain, _ := c.Params.Get("domain")

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return nil
	}

	database := db.Database{}
	database.Connect(dbFile)
	return &database
}

func GetSummaries(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	sessions := database.GetSessions(c)
	pageViews := database.GetPageViews(c)
	avgSessionDuration := database.GetAvgSessionDuration(c)
	bounceRate := database.GetBounceRate(c)

	c.IndentedJSON(http.StatusOK, &analytics.SummaryResponse{
		Sessions:           sessions,
		PageViews:          pageViews,
		AvgSessionDuration: avgSessionDuration,
		BounceRate:         bounceRate,
	})
}

func GetBrowsers(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetBrowsers(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
		return
	}

	items := arrayFromRows(rows)

	previousFilters := make([]string, 0)

	browser, hasBrowser := c.GetQuery("b")
	browserVersion, hasBrowserVersion := c.GetQuery("bv")

	if hasBrowser {
		previousFilters = append(previousFilters, browser)
	}

	if hasBrowserVersion {
		bver := strings.Split(browserVersion, "/")
		previousFilters = append(previousFilters, bver...)
	}

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		PreviousFilters: previousFilters,
		Items:           items,
	})
}

func GetOSs(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetOSs(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get OSs")
		return
	}

	items := arrayFromRows(rows)

	previousFilters := make([]string, 0)

	os, hasOs := c.GetQuery("os")
	osVersion, hasOsVersion := c.GetQuery("osv")

	if hasOs {
		previousFilters = append(previousFilters, os)
	}

	if hasOsVersion {
		osver := strings.Split(osVersion, "/")
		previousFilters = append(previousFilters, osver...)
	}

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		PreviousFilters: previousFilters,
		Items:           items,
	})
}

func GetCountries(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetCountries(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Countries")
		return
	}

	items := arrayFromRows(rows)

	country, hasCountry := c.GetQuery("c")

	previousFilters := make([]string, 0)

	if hasCountry {
		previousFilters = append(previousFilters, country)
	}

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		PreviousFilters: previousFilters,
		Items:           items,
	})
}

func GetReferrers(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetReferrers(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Referrers")
		return
	}

	items := arrayFromRows(rows)

	referrer, hasReferrer := c.GetQuery("r")

	previousFilters := make([]string, 0)

	if hasReferrer {
		previousFilters = append(previousFilters, referrer)
	}

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		PreviousFilters: previousFilters,
		Items:           items,
	})
}

func GetPages(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetPages(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Pages")
		return
	}

	items := arrayFromRows(rows)

	path, hasPath := c.GetQuery("pg")

	previousFilters := make([]string, 0)

	if hasPath {
		previousFilters = append(previousFilters, path)
	}

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		PreviousFilters: previousFilters,
		Items:           items,
	})
}

func GetWebsites(c *gin.Context) {
	sites := config.Config.Websites
	c.IndentedJSON(http.StatusOK, &sites)
}
