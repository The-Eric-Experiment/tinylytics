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
