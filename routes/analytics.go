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

	c.IndentedJSON(http.StatusOK, &analytics.VersionedListResponse{
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

	c.IndentedJSON(http.StatusOK, &analytics.VersionedListResponse{
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

	c.IndentedJSON(http.StatusOK, &analytics.VersionedListResponse{
		Items: items,
	})
}

func arrayFromRows(rows *sql.Rows, database *db.Database) []*analytics.Versioned {
	list := make([]*analytics.Versioned, 0)
	for rows.Next() {
		var output analytics.Versioned

		database.Scan(rows, &output)

		list = append(list, &output)
	}
	return list
}
