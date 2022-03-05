package helpers

import (
	"errors"
	"log"
	"path"
	conf "tinylytics/config"

	"github.com/google/uuid"
)

func FindWebsite(domain string) (*conf.WebsiteConfig, error) {
	for _, element := range conf.Config.Websites {
		if domain == element.Domain {
			return &element, nil
		}
	}
	return nil, errors.New("Site not found")
}

func GetDatabaseFileName(domain string) string {
	site, err := FindWebsite(domain)
	if err != nil {
		log.Fatal(err)
	}

	nm, _ := uuid.FromBytes([]byte("0d032761-6264-49d4-b099-74219d6d564d"))
	dbHash := uuid.NewSHA1(nm, []byte(site.Domain)).String()

	filename := dbHash + ".db"
	return path.Join(conf.Config.DataFolder, filename)
}
