package ua

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tinylytics/constants"
	"tinylytics/helpers"

	"github.com/ua-parser/uap-go/uaparser"
)

type UA struct {
	Browser      string `json:"browser"`
	BrowserMajor string `json:"browserMajor"`
	BrowserMinor string `json:"browserMinor"`
	BrowserPatch string `json:"browserPatch"`
	OS           string `json:"os"`
	OSMajor      string `json:"osMajor"`
	OSMinor      string `json:"osMinor"`
	OSPatch      string `json:"osPatch"`
}

func ParseUA(uagent string) UA {
	parser, err := uaparser.New(helpers.GetDataPath(constants.UA_REGEX_FILE_NAME))
	if err != nil {
		log.Fatal(err)
	}

	client := parser.Parse(uagent)

	return UA{
		Browser:      client.UserAgent.Family,
		BrowserMajor: client.UserAgent.Major,
		BrowserMinor: client.UserAgent.Minor,
		BrowserPatch: client.UserAgent.Patch,
		OS:           client.Os.Family,
		OSMajor:      client.Os.Major,
		OSMinor:      client.Os.Minor,
		OSPatch:      client.Os.Patch,
	}
}

func downloadUARegex(filepath string) error {
	// Get the data
	resp, err := http.Get(constants.UA_DOWNLOAD_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Initialize() {
	e, _ := helpers.Exists(helpers.GetDataPath(constants.UA_REGEX_FILE_NAME))
	if e {
		fmt.Println("Aready exists, UA Regex")
		return
	}

	fmt.Println("Downloading UA Regex")

	err := downloadUARegex(helpers.GetDataPath(constants.UA_REGEX_FILE_NAME))

	if err != nil {
		panic(err)
	}

	fmt.Println("UA Regex downloaded")
}
