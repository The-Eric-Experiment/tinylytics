package ua

import (
	"log"

	"github.com/ua-parser/uap-go/uaparser"
)

type UA struct {
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browserVersion"`
	OS             string `json:"os"`
	OSVersion      string `json:"osVersion"`
}

func ParseUA(uagent string) UA {
	parser, err := uaparser.New("./regexes.yaml")
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
