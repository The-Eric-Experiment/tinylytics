package main

import (
	"fmt"
	"html/template"
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

	// Load HTML templates with custom functions
	router.SetFuncMap(template.FuncMap{
		"jsonItems": routes.JSONItems,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: number of arguments must be even")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"seq": func(n int) []int {
			result := make([]int, n)
			for i := range result {
				result[i] = i
			}
			return result
		},
	})
	router.LoadHTMLGlob("templates/*.html")

	// Serve static files (CSS, JS, images)
	router.Static("/static", "./static")

	// API routes for event tracking
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

	// HTML template routes using query params to avoid greedy route matching
	router.GET("/", routes.GetAnalyticsPage)
	router.GET("/analytics", routes.GetAnalyticsPage)

	// HTML fragment routes for HTMX (using same endpoints as API but with Accept header or query param)
	router.GET("/summary-cards", routes.GetSummaries)
	router.GET("/browsers-table", routes.GetBrowsers)
	router.GET("/os-table", routes.GetOSs)
	router.GET("/pages-table", routes.GetPages)
	router.GET("/referrers-table", routes.GetReferrers)
	router.GET("/countries-table", routes.GetCountries)

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
