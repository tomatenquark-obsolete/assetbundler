package main

import "C"

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"github.com/shibukawa/configdir"
	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
	"github.com/tomatenquark/assetbundler/internal/resources"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Defines the status of the download
const (
	CONNECTION_ABORTED = -1
	IN_PROGRESS = iota
	FINISHED = iota
)

var downloads = make(map[string]int)

// Downloads a batch of files.
// Returns a slice of destination paths
func DownloadResourcesAndArchive(sources []url.URL, destinations []string, zipPath string, serverDirectory string) {
	requests := make([]*grab.Request, 0)
	client := grab.NewClient()

	for index, source := range sources {
		request, _ := grab.NewRequest(destinations[index], source.String())
		requests = append(requests, request)
	}

	responsesChannel := client.DoBatch(5, requests...)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for resp := range responsesChannel {
		/*if err := resp.Err(); err != nil {
			panic(err)
		}*/
		fmt.Printf("Downloaded %s to %s\n", resp.Request.URL(), resp.Filename)
	}

	err := archive.ZipFiles(zipPath, destinations, serverDirectory)
	if err != nil {
		downloads[zipPath] = CONNECTION_ABORTED
	}
	downloads[zipPath] = FINISHED
}

// Starts a download for the given map from servercontent.
// Returns the destination of the resulting archive.
//export StartDownload
func StartDownload(servercontent *C.char, servermap *C.char) *C.char {
	// Create a temporary archive for the download
	tempFile, _ := ioutil.TempFile("", "maparchive*.zip")
	tempFile.Close()
	os.Remove(tempFile.Name())

	// Verify that source is indeed a URL
	serverContent := C.GoString(servercontent)
	mapString := C.GoString(servermap)
	configPath := fmt.Sprint("/packages/base/", mapString, ".cfg")
	uri, err := url.Parse(fmt.Sprint(serverContent, configPath))
	if err != nil {
		return C.CString("")
	}
	// Use ~/tomatenquark/packages/servername to store packages
	configDirectories := configdir.New("tomatenquark", "")
	userDirectories := configDirectories.QueryCacheFolder()
	hostIndex := strings.Index(uri.String(), uri.Hostname())
	configIndex := strings.Index(uri.String(), configPath)
	serverDirectory := path.Join(userDirectories.Path, uri.String()[hostIndex:configIndex])

	// Gather all the resources from the map config file
	res, err := resources.Collect(*uri, serverDirectory)
	if err != nil {
		downloads[tempFile.Name()] = CONNECTION_ABORTED
		return C.CString("")
	}

	// Add additional map resources
	mapFiles := []string{"ogz", "wpt", "jpg"}
	for _, mapFile := range mapFiles {
		res = append(res, config.Resource{"map", path.Join("base", strings.Replace(path.Base(uri.Path), "cfg", mapFile, 1))})
	}

	// Prepare download list
	var sources []url.URL
	var destinations []string
	for _, resource := range res {
		resourceURI := *uri
		var resourcePath string
		switch resource.Property {
		case "mapsound":
			resourcePath = path.Join("sounds", resource.Path)
		default:
			resourcePath = resource.Path
		}
		resourceURI.Path = strings.Replace(resourceURI.Path, configPath, fmt.Sprint("/packages/", resourcePath), 1)
		sources = append(sources, resourceURI)
		destinations = append(destinations, path.Join(serverDirectory, path.Join("packages", resourcePath)))
	}

	// Dispatch download
	destinations = append(destinations, path.Join(serverDirectory, uri.Path))
	downloads[tempFile.Name()] = IN_PROGRESS
	go DownloadResourcesAndArchive(sources, destinations, tempFile.Name(), serverDirectory)
	return C.CString(tempFile.Name())
}

// Returns the status of the download
//export GetStatus
func GetStatus(zipPath *C.char) C.int {
	safeZipPath := C.GoString(zipPath)
	return C.int(downloads[safeZipPath])
}

func main() {}
