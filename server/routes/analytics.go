package routes

import (
	"database/sql"
	"net/http"
	"tinylytics/analytics"
	"tinylytics/db"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
)

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

	c.IndentedJSON(http.StatusOK, &analytics.SummaryResponse{
		Sessions:  sessions,
		PageViews: pageViews,
	})
}

func GetBrowsers(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetBrowsers(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
	}

	items := arrayFromRows(rows, database)

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		Items: items,
	})
}

func GetOSs(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetOSs(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get OSs")
	}

	items := arrayFromRows(rows, database)

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		Items: items,
	})
}

func GetCountries(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetCountries(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Countries")
	}

	items := arrayFromRows(rows, database)

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		Items: items,
	})
}

func GetReferrers(c *gin.Context) {
	database := getDB(c)
	defer database.Close()

	rows, err := database.GetReferrers(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get Referrers")
	}

	items := arrayFromRows(rows, database)

	c.IndentedJSON(http.StatusOK, &analytics.AnalyticsListResponse{
		Items: items,
	})
}

func arrayFromRows(rows *sql.Rows, database *db.Database) []*analytics.AnalyticsItem {
	list := make([]*analytics.AnalyticsItem, 0)
	for rows.Next() {
		var output analytics.AnalyticsItem

		database.Scan(rows, &output)

		list = append(list, &output)
	}
	return list
}
