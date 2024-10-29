package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"tinylytics/config"
	conf "tinylytics/config"
	"tinylytics/db"
	"tinylytics/event"
	"tinylytics/geo"
	"tinylytics/helpers"
	"tinylytics/routes"
	"tinylytics/ua"

	"github.com/gin-gonic/contrib/static"
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
		filename, err := helpers.GetDatabaseFileName(element.Domain)

		if err != nil {
			panic(err)
		}

		fmt.Println("Initializing database for domain: ", element.Domain)

		initializeDb(filename)
	}
}

func init() {
	conf.LoadConfig("config.yaml")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	newpath := filepath.Join(wd, config.Config.DataFolder)
	err = os.MkdirAll(newpath, os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}

	ua.Initialize()
	geo.Initialize()

	initializeDatabases()

	eventQueue.Connect()
}

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/event", routes.PostEvent(&eventQueue))
		api.GET("/sites", routes.GetWebsites)
		api.GET("/:domain/summaries", routes.GetSummaries)
		api.GET("/:domain/browsers", routes.GetBrowsers)
		api.GET("/:domain/os", routes.GetOSs)
		api.GET("/:domain/countries", routes.GetCountries)
		api.GET("/:domain/pages", routes.GetPages)
		api.GET("/:domain/referrers", routes.GetReferrers)
	}

	// I don't like this, maybe change to echo
	router.Use(static.Serve("/", static.LocalFile("./client", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("./client/index.html")
	})

	eventQueue.Listen(event.ProcessEvent)

	router.Run("0.0.0.0:8099")
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
