package routes

import (
	"net/http"
	"tinylytics/analytics"
	"tinylytics/db"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
)

func GetSummaries(c *gin.Context) {
	domain, _ := c.GetQuery("d")

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return
	}

	database := db.Database{}
	database.Connect(dbFile)
	defer database.Close()

	sessions := database.GetSessions()
	pageViews := database.GetPageViews()

	c.IndentedJSON(http.StatusOK, &analytics.SummaryResponse{
		Sessions:  sessions,
		PageViews: pageViews,
	})
}

func GetBrowsers(c *gin.Context) {
	domain, _ := c.GetQuery("d")

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return
	}

	database := db.Database{}
	database.Connect(dbFile)
	defer database.Close()

	rows, err := database.GetBrowsers()

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
	}

	browserList := make([]*analytics.Browser, 0)

	for rows.Next() {
		var name string
		var count int64
		rows.Scan(&name, &count)

		browserList = append(browserList, &analytics.Browser{
			Name:  name,
			Count: count,
		})
	}

	c.IndentedJSON(http.StatusOK, &analytics.BrowserListResponse{
		Items: browserList,
	})
}
