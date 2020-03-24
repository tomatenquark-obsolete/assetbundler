package main

import (
	"bufio"
	"fmt"
	"github.com/shibukawa/configdir"
	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)


// Downloads a config, filters it and writes it to the given destination.
// Returns a list of resources used in the config file.
func ParseConfig(source url.URL, destination string) ([]config.Resource, error) {
	tempFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return nil, err
	}
	tempFile.Close()
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
	config, err := os.Open(destination)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(config)
	reader.WriteTo(writer)
	return resources, nil
}

// Takes a list of resources and returns the path of all exec type resources
func FilterResourcesForExec(resources []config.Resource) []string {
	var configs []string

	for _, resource := range resources {
		if resource.Property == "exec" {
			configs = append(configs, resource.Path)
		}
	}
	return configs
}


// Downloads a map from the given path to disk cache and returns a
// path to a ZIP archive packaged with all the necessary contents.
//export DownloadMap
func DownloadMap(source string, progress func(int, int)) (string, error) {
	// Verify that source is indeed a URL
	uri, err := url.Parse(source)
	if err != nil {
		return "", err
	}
	// Use ~/tomatenquark/packages/servername to store packages
	userDirectories := configdir.New("tomatenquark", "packages")
	localDirectory := userDirectories.LocalPath
	serverDirectory := path.Join(localDirectory, uri.Hostname())

	// Gather all the resources from the map config file
	var resources []config.Resource
	var configResources []config.Resource
	configResources, err = ParseConfig(*uri, path.Join(serverDirectory, path.Base(uri.Path)))
	if err != nil {
		return "", err
	}
	resources = append(resources, configResources...)

	// Repeat this for all the configs until all resources have been aggregated
	configs := make([]string, 0)
	configs = append(configs, FilterResourcesForExec(resources)...)

	for len(configs) > 0 {
		configPath, configs := configs[len(configs)-1], configs[:len(configs)-1]
		configURI := *uri
		configURI.Path = configPath
		configResources, err = ParseConfig(configURI, path.Join(serverDirectory, configPath))
		resources = append(resources, configResources...)
		configs = append(configs, FilterResourcesForExec(resources)...)
	}

	// Start downloading

	// Package all the destination files into a single ZIP

	// Return the path of the zip
	return "", nil
}

func main() {
	// Take path as input
	// Download map.cfg
	// Gather all the resources
	// Gather all cfg files
	//
	//
	fmt.Printf("ohmygod")
}