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

	browser, hasBrowser := c.GetQuery("b")
	browserMajor, hasBrowserMajor := c.GetQuery("bm")

	dbFile, err := helpers.GetDatabaseFileName(domain)

	if err != nil {
		c.String(http.StatusBadRequest, "", err)
		return
	}

	database := db.Database{}
	database.Connect(dbFile)
	defer database.Close()

	var b *string = nil
	if hasBrowser {
		b = &browser
	}

	var bv *string = nil
	if hasBrowserMajor {
		bv = &browserMajor
	}

	rows, err := database.GetBrowsersQuery(&db.QueryFilters{
		Browser:      b,
		BrowserMajor: bv,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Couldn't get browsers")
	}

	browserList := make([]*analytics.Browser, 0)

	for rows.Next() {
		var name string
		var browserMajor int64
		var browserMinor int64
		var count int64

		if hasBrowser {
			rows.Scan(&name, &count, &browserMajor)

			if hasBrowserMajor {
				rows.Scan(&name, &count, &browserMajor, &browserMinor)
			}
		} else {
			rows.Scan(&name, &count)
		}

		browserList = append(browserList, &analytics.Browser{
			Name:         name,
			Count:        count,
			BrowserMajor: browserMajor,
			BrowserMinor: browserMinor,
		})
	}

	c.IndentedJSON(http.StatusOK, &analytics.BrowserListResponse{
		Items: browserList,
	})
}
