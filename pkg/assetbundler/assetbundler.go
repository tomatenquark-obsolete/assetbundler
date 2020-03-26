package main

import "C"

import (
	"bufio"
	"github.com/shibukawa/configdir"
	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

// Returns the index of a given resource
func Index(vs []config.Resource, t config.Resource) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Downloads a config, filters it and writes it to the given destination.
// Returns a list of resources used in the config file.
func ParseConfig(source url.URL, destination string) ([]config.Resource, error) {
	tempFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return nil, err
	}
	tempFile.Close()
	os.Remove(tempFile.Name())
	tempPath, err := archive.DownloadFile(source, tempFile.Name())
	if err != nil {
		return nil, err
	}
	tempConfig, err := os.Open(tempPath)
	if err != nil {
		return nil, err
	}
	defer tempConfig.Close()
	reader := bufio.NewReader(tempConfig)
	filteredConfig, err := config.Filter(reader)
	if err != nil {
		return nil, err
	}
	resources, err := config.ReadResources(filteredConfig)
	if err != nil {
		return nil, err
	}
	os.MkdirAll(path.Dir(destination), os.ModePerm)
	config, err := os.Create(destination)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(config)
	reader.WriteTo(writer)
	return resources, nil
}

// Takes a list of resources and returns all exec type resources
func GetConfigs(resources []config.Resource) []config.Resource {
	var configs []config.Resource

	for _, resource := range resources {
		if resource.Property == "exec" {
			configs = append(configs, resource)
		}
	}
	return configs
}

// Downloads all the configs starting from "start" and returns a list of all the collected resources
func CollectResources(start url.URL, destinationDirectory string) ([]config.Resource, error) {
	var resources []config.Resource
	var configResources []config.Resource
	configResources, err := ParseConfig(start, path.Join(destinationDirectory, start.Path))
	if err != nil {
		return nil, err
	}
	resources = append(resources, configResources...)

	configs := make([]config.Resource, 0)
	configs = append(configs, GetConfigs(resources)...)

	// Repeat for all the configs until all resources have been aggregated
	for len(configs) > 0 {
		// Pop element from configs
		conf, configs := configs[len(configs)-1], configs[:len(configs)-1]
		// Pop element from resources
		resourceIndex := Index(resources, conf)
		resources := append(resources[:resourceIndex], resources[resourceIndex+1:]...)

		// Download config and evaluate it
		configURI := start
		configURI.Path = conf.Path
		configResources, err = ParseConfig(configURI, path.Join(destinationDirectory, conf.Path))
		resources = append(resources, configResources...)
		configs = append(configs, GetConfigs(resources)...)
	}

	return resources, nil
}

// Downloads a map from the given path to disk cache and returns a
// path to a ZIP archive packaged with all the necessary contents.
//export DownloadMap
func DownloadMap(servercontent *C.char, servermap *C.char) *C.char {
	// Verify that source is indeed a URL
	serverContent := C.GoString(servercontent)
	mapString := C.GoString(servermap)
	uri, err := url.Parse(path.Join(serverContent, "packages", mapString, ".cfg"))
	if err != nil {
		return C.CString("")
	}
	// Use ~/tomatenquark/packages/servername to store packages
	configDirectories := configdir.New("tomatenquark", "")
	userDirectories := configDirectories.QueryCacheFolder()
	serverDirectory := path.Join(userDirectories.Path, uri.Hostname())
	// Gather all the resources from the map config file
	resources, err := CollectResources(*uri, serverDirectory)
	// Also add the map and waypoint as a download resources
	resources = append(resources, config.Resource{"map", path.Join("base", strings.Replace(path.Base(uri.Path), "cfg", "ogz", 1))})
	resources = append(resources, config.Resource{"map", path.Join("base", strings.Replace(path.Base(uri.Path), "cfg", "wpt", 1))})
	if err != nil {
		return C.CString("")
	}

	// Start downloading
	var sources []url.URL
	var destinations []string
	for _, resource := range resources {
		resourceURI := *uri
		var resourcePath string
		switch resource.Property {
		case "mmodel":
			resourcePath = path.Join("models", resource.Path)
		case "mapsound":
			resourcePath = path.Join("sounds", resource.Path)
		default:
			resourcePath = resource.Path
		}
		resourceURI.Path = path.Join("packages", resourcePath)
		sources = append(sources, resourceURI)
		destinations = append(destinations, path.Join(serverDirectory, path.Join("packages", resourcePath)))
	}

	_, err = archive.DownloadBatch(sources, destinations)
	if err != nil {
		return C.CString("")
	}

	// Package all the destination files into a single ZIP
	tempFile, err := ioutil.TempFile("", "maparchive*.zip")
	if err != nil {
		return C.CString("")
	}
	tempFile.Close()
	os.Remove(tempFile.Name())
	destinations = append(destinations, path.Join(serverDirectory, uri.Path))
	if err := archive.ZipFiles(tempFile.Name(), destinations, serverDirectory); err != nil {
		return C.CString("")
	}

	// Return the path of the zip
	return C.CString(tempFile.Name())
}

func main() {
	argsWithoutProg := os.Args[1:]
	host, serverMap := argsWithoutProg[0], argsWithoutProg[1]
	safeHost := C.CString(host)
	safeServerMap := C.CString(serverMap)
	DownloadMap(safeHost, safeServerMap)
}
