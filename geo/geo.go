package geo

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/oschwald/geoip2-golang"
)

var url = "https://raw.githubusercontent.com/GitSquared/node-geolite2-redist/master/redist/GeoLite2-Country.tar.gz"
var destinationFile = "data/GeoLite2-Country.tar.gz"

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func extractTarGz(gzipStream io.Reader) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		nm := header.Name
		fmt.Println(path.Ext(nm))
		if header.Typeflag == tar.TypeReg && path.Ext(nm) == ".mmdb" {
			file := path.Base(nm)
			destFile := path.Join("data", file)

			outFile, err := os.Create(destFile)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()
		}

		// 	switch header.Typeflag {
		// 	case tar.TypeDir:
		// 		if err := os.Mkdir(header.Name, 0755); err != nil {
		// 			log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
		// 		}
		// 	case tar.TypeReg:
		// 		nm := header.Name

		// 		fmt.Println(nm)
		// 		outFile, err := os.Create(nm)
		// 		if err != nil {
		// 			log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
		// 		}
		// 		if _, err := io.Copy(outFile, tarReader); err != nil {
		// 			log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
		// 		}
		// 		outFile.Close()

		// 	default:
		// 		log.Fatalf(
		// 			"ExtractTarGz: uknown type: %s in %s",
		// 			header.Typeflag,
		// 			header.Name)
		// 	}

	}
}

func downloadGeoLite(filepath string) error {

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
	e, _ := exists(destinationFile)
	if e {
		fmt.Println("Aready exists, geolite")
		return
	}

	fmt.Println("Downloading Geolite")

	err := downloadGeoLite(destinationFile)

	if err != nil {
		panic(err)
	}

	fmt.Println("Extracting Geolite")

	r, err := os.Open(destinationFile)
	if err != nil {
		fmt.Println("error")
	}
	extractTarGz(r)

	fmt.Println("GeoliteDownloaded")
}

func GetGeo(ipInput string) string {
	db, err := geoip2.Open("data/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipInput)
	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	return record.Country.IsoCode
	// fmt.Printf("Portuguese (BR) city name: %v\n", record.Country.Names["en-US"])
	// // fmt.Printf("Confidence", record.Country.Confidence)
	// fmt.Printf("IsoCode", record.Country.IsoCode)
	// fmt.Printf("GeoNameID", record.Country.GeoNameID)
	// // fmt.Printf("IsEuropeanUnion", record.Country.IsEuropeanUnion)
}
