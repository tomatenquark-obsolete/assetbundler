package resources

import (
	"bufio"
	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

// Returns the index of a given resource
func Index(vs []config.Resource, t config.Resource) int {
	for i, v := range vs {
		if v.Path == t.Path {
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
func Collect(start url.URL, destinationDirectory string) ([]config.Resource, error) {
	var resources []config.Resource
	var configResources []config.Resource
	var err error
	configResources, err = ParseConfig(start, path.Join(destinationDirectory, start.Path))
	if err != nil {
		return nil, err
	}
	resources = append(resources, configResources...)

	configs, processed := make([]config.Resource, 0), make([]config.Resource, 0)
	configs = append(configs, GetConfigs(resources)...)

	// Repeat for all the configs until all resources have been aggregated
	for len(configs) > 0 {
		// Pop element from configs
		var conf config.Resource
		conf, configs = configs[len(configs)-1], configs[:len(configs)-1]
		// Never process a file twice
		confIndex := Index(processed, conf)
		if confIndex >= 0 {
			continue;
		}
		processed = append(processed, conf)
		// Pop element from resources
		resourceIndex := Index(resources, conf)
		resources = append(resources[:resourceIndex], resources[resourceIndex+1:]...)

		// Download config and evaluate it
		configURI := start
		configURI.Path = conf.Path
		configResources, err = ParseConfig(configURI, path.Join(destinationDirectory, conf.Path))
		resources = append(resources, configResources...)
		configs = append(configs, GetConfigs(resources)...)
	}

	return resources, nil
}
