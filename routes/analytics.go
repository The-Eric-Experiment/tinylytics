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

	dbFile := helpers.GetDatabaseFileName(domain)

	database := db.Database{}
	database.Connect(dbFile)
	defer database.Close()

	sessions := database.GetSessions()

	c.IndentedJSON(http.StatusOK, &analytics.SummaryResponse{
		Sessions: sessions,
	})
}
