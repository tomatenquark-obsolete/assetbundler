package resources

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/tomatenquark/assetbundler/internal/archive"
	"github.com/tomatenquark/assetbundler/internal/config"
)

// Returns the index of a given resource
func Index(vs []config.FileToLoad, t config.FileToLoad) int {
	for i, v := range vs {
		if v.Path == t.Path {
			return i
		}
	}
	return -1
}

// Downloads a config, filters it and writes it to the given destination.
// Returns a list of resources used in the config file.
func ParseConfig(source url.URL, destination string) ([]config.FileToLoad, error) {
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
	tempConfig, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return nil, err
	}
	filteredConfig, err := config.Filter(string(tempConfig))
	if err != nil {
		return nil, err
	}
	resources, err := config.FilesToLoadFromCfg(strings.NewReader(filteredConfig))
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Dir(destination), 0644)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(destination, []byte(filteredConfig), 0644)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// Takes a list of resources and returns all exec type resources
func GetConfigs(resources []config.FileToLoad) []config.FileToLoad {
	var configs []config.FileToLoad

	for _, resource := range resources {
		if resource.Command == "exec" {
			configs = append(configs, resource)
		}
	}
	return configs
}

// Downloads all the configs starting from "start" and returns a list of all the collected resources
func Collect(start url.URL, destinationDirectory string) ([]config.FileToLoad, error) {
	var resources []config.FileToLoad
	var configResources []config.FileToLoad
	var err error
	configPath := fmt.Sprint("packages/base/", path.Base(start.Path))
	configResources, err = ParseConfig(start, path.Join(destinationDirectory, configPath))
	if err != nil {
		return nil, err
	}
	resources = append(resources, configResources...)

	configs, processed := make([]config.FileToLoad, 0), make([]config.FileToLoad, 0)
	configs = append(configs, GetConfigs(resources)...)

	// Repeat for all the configs until all resources have been aggregated
	for len(configs) > 0 {
		// Pop element from configs
		var conf config.FileToLoad
		conf, configs = configs[len(configs)-1], configs[:len(configs)-1]
		// Never process a file twice
		confIndex := Index(processed, conf)
		if confIndex >= 0 {
			continue
		}
		processed = append(processed, conf)
		// Pop element from resources
		resourceIndex := Index(resources, conf)
		resources = append(resources[:resourceIndex], resources[resourceIndex+1:]...)

		// Download config and evaluate it
		configURI := start
		configURI.Path = strings.Replace(configURI.Path, configPath, conf.Path, 1)
		confPath := path.Join(destinationDirectory, conf.Path)
		configResources, err = ParseConfig(configURI, confPath)
		resources = append(resources, configResources...)
		configs = append(configs, GetConfigs(resources)...)
	}

	return resources, nil
}
