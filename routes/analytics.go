package routes

import (
	"fmt"
	"net/http"
	"tinylytics/analytics"
	"tinylytics/db"
	"tinylytics/helpers"

	"github.com/gin-gonic/gin"
)

func GetSummaries(c *gin.Context) {
	domain, _ := c.Params.Get("domain")

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
	domain, _ := c.Params.Get("domain")

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return
	}

	database := db.Database{}
	database.Connect(dbFile)
	defer database.Close()

	rows, err := database.GetBrowsers(c)

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
	}

	browserList := make([]*analytics.Browser, 0)

	for rows.Next() {
		var output analytics.Browser

		database.Scan(rows, &output)

		fmt.Println(output)

		browserList = append(browserList, &output)
	}

	c.IndentedJSON(http.StatusOK, &analytics.BrowserListResponse{
		Items: browserList,
	})
}
