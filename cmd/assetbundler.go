package main

import "C"
import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/shibukawa/configdir"

	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
	"github.com/tomatenquark/assetbundler/internal/resources"
)

// Downloads a map from the given path to disk cache and returns a
// path to a ZIP archive packaged with all the necessary contents.
func DownloadMap(servercontent string, servermap string) string {
	// Verify that source is indeed a URL
	configPath := fmt.Sprint("/packages/base/", servermap, ".cfg")
	uri, err := url.Parse(fmt.Sprint(servercontent, configPath))
	if err != nil {
		return ""
	}
	// Use ~/tomatenquark/packages/servername to store packages
	configDirectories := configdir.New("tomatenquark", "")
	userDirectories := configDirectories.QueryCacheFolder()
	hostIndex := strings.Index(uri.String(), uri.Hostname())
	configIndex := strings.Index(uri.String(), configPath)
	serverDirectory := path.Join(userDirectories.Path, uri.String()[hostIndex:configIndex])
	// Gather all the resources from the map config file
	resources, err := resources.Collect(*uri, serverDirectory)
	// Also add the map and waypoint as a download resources
	mapFiles := []string{"ogz", "wpt", "jpg"}
	for _, mapFile := range mapFiles {
		resources = append(resources, config.FileToLoad{"map", path.Join("base", strings.Replace(path.Base(uri.Path), "cfg", mapFile, 1))})
	}
	if err != nil {
		return ""
	}

	// Start downloading
	var sources []url.URL
	var destinations []string
	for _, resource := range resources {
		resourceURI := *uri
		var resourcePath string
		switch resource.Command {
		case "mapsound":
			resourcePath = path.Join("sounds", resource.Path)
		default:
			resourcePath = resource.Path
		}
		resourceURI.Path = strings.Replace(resourceURI.Path, configPath, fmt.Sprint("/packages/", resourcePath), 1)
		sources = append(sources, resourceURI)
		destinations = append(destinations, path.Join(serverDirectory, path.Join("packages/", resourcePath)))
	}

	_, err = archive.DownloadBatch(sources, destinations)
	if err != nil {
		return ""
	}

	// Package all the destination files into a single ZIP
	tempFile, err := ioutil.TempFile("", "maparchive*.zip")
	if err != nil {
		return ""
	}
	tempFile.Close()
	os.Remove(tempFile.Name())
	destinations = append(destinations, path.Join(serverDirectory, configPath))
	if err := archive.ZipFiles(tempFile.Name(), destinations, serverDirectory); err != nil {
		return ""
	}

	// Return the path of the zip
	return tempFile.Name()
}

func main() {
	argsWithoutProg := os.Args[1:]
	host, serverMap := argsWithoutProg[0], argsWithoutProg[1]
	archivePath := DownloadMap(host, serverMap)
	fmt.Println(archivePath)
}
