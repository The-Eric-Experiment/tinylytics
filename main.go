package main

import (
	"fmt"
	"net/http"
	"runtime"
	conf "tinylytics/config"
	"tinylytics/db"
	"tinylytics/event"
	"tinylytics/geo"
	"tinylytics/helpers"
	"tinylytics/routes"
	"tinylytics/ua"

	"github.com/gin-gonic/gin"
)

var eventQueue = event.EventQueue{}

func initializeDb(filename string) {
	database := db.Database{}
	database.Connect(filename)
	defer database.Close()
	database.Initialize()
}

func initializeDatabases() {
	for _, element := range conf.Config.Websites {
		filename := helpers.GetDatabaseFileName(element.Domain)
		initializeDb(filename)
	}
}

func init() {
	conf.LoadConfig("config.yaml")

	ua.Initialize()
	geo.Initialize()

	initializeDatabases()

	eventQueue.Connect()
}

func main() {
	router := gin.Default()
	router.POST("/api/event", routes.PostEvent(&eventQueue))
	router.GET("/ua", func(c *gin.Context) {
		uagent := c.Request.Header.Get("User-Agent")
		result := ua.ParseUA(uagent)
		router.LoadHTMLGlob("*.html")
		c.HTML(http.StatusOK, "ua.html", gin.H{
			"UA":             uagent,
			"Browser":        result.Browser,
			"BrowserVersion": result.BrowserVersion,
			"OS":             result.OS,
			"OSVersion":      result.OSVersion,
		})
	})

	eventQueue.Listen(event.ProcessEvent)

	router.Run("localhost:8008")
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
