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
	router.POST("/api/event", routes.PostEvent(&eventQueue))
	router.GET("/analytics/summaries", routes.GetSummaries)
	router.GET("/analytics/browsers", routes.GetBrowsers)

	eventQueue.Listen(event.ProcessEvent)

	router.Run()
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
