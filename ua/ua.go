package ua

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tinylytics/helpers"

	"github.com/ua-parser/uap-go/uaparser"
)

var url = "https://raw.githubusercontent.com/ericmackrodt/uap-core/master/regexes.yaml"
var destinationFile = "data/regexes.yaml"

type UA struct {
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browserVersion"`
	OS             string `json:"os"`
	OSVersion      string `json:"osVersion"`
}

func ParseUA(uagent string) UA {
	parser, err := uaparser.New("./data/regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client := parser.Parse(uagent)

	return UA{
		Browser:        client.UserAgent.Family,
		BrowserVersion: client.UserAgent.Major,
		OS:             client.Os.Family,
		OSVersion:      client.Os.Major,
	}
}

func downloadUARegex(filepath string) error {
	// Get the data
	resp, err := http.Get(url)
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
	e, _ := helpers.Exists(destinationFile)
	if e {
		fmt.Println("Aready exists, UA Regex")
		return
	}

	fmt.Println("Downloading UA Regex")

	err := downloadUARegex(destinationFile)

	if err != nil {
		panic(err)
	}

	fmt.Println("UA Regex downloaded")
}
