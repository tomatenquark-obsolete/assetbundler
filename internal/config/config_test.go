package config

import (
	"bufio"
	"fmt"
	"github.com/tomatenquark/assetbundler/internal/util"
	"os"
	"path"
	"testing"
)

func TestReadResources(t *testing.T) {
	projectDirectory := util.GetProjectDirectory()
	testData := path.Join(projectDirectory, "internal/config/testdata")
	configPath := path.Join(testData, "collide.cfg")
	configFile, err := os.Open(configPath)
	if err != nil {
		t.Error(err)
	}
	defer configFile.Close()
	reader := bufio.NewReader(configFile)
	resources, err := ReadResources(reader)
	if err != nil {
		t.Error(err)
	}

	if len(resources) != 46 {
		t.Error("There should be 48 resources.")
	}
}

func TestFilter(t *testing.T) {
	projectDirectory := util.GetProjectDirectory()
	testData := path.Join(projectDirectory, "internal/config/testdata")
	configPath := path.Join(testData, "corrupted.cfg")
	configFile, err := os.Open(configPath)
	if err != nil {
		t.Error(err)
	}
	defer configFile.Close()
	reader := bufio.NewReader(configFile)
	output, err := Filter(reader)
	if err != nil {
		t.Error(err)
	}

	count := 0
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		t.Error(err)
	}

	if count != 8 {
		t.Log(count)
		t.Error("Only 8 lines should remain")
	}
}

func TestComplete(t *testing.T) {
	projectDirectory := util.GetProjectDirectory()
	testData := path.Join(projectDirectory, "internal/config/testdata")
	configPath := path.Join(testData, "corrupted.cfg")
	configFile, err := os.Open(configPath)
	if err != nil {
		t.Error(err)
	}
	defer configFile.Close()

	reader := bufio.NewReader(configFile)
	output, err := Filter(reader)
	if err != nil {
		t.Error(err)
	}

	resources, err := ReadResources(output)
	if err != nil {
		if err.Error() != fmt.Sprintln("To few arguments in line 3") {
			t.Error("Invalid error")
		}
	}

	if len(resources) != 0 {
		t.Log(resources)
		t.Error("Invalid amount of resources")
	}
}
