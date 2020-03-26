package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/tomatenquark/assetbundler/internal/collection"
	"io"
	"regexp"
	"strings"
)

// Maps the relation between ICOMMAND and the parameters it takes (file path)
var Properties = map[string]int{
	"texture": 2,
	"mmodelfile": 1,
	"mapsound": 1,
	"skybox": 1,
	"exec": 1,
	"cloudlayer": 1,
}

// Whitelist contains all ICOMMAND instructions which are valid in map configs
var Whitelist = []string{
	"ambient",
	"autograss",
	"base_1",
	"base_10",
	"base_2",
	"base_3",
	"base_4",
	"base_5",
	"base_6",
	"base_7",
	"base_8",
	"base_9",
	"blurlms",
	"blurskylight",
	"causticmillis",
	"causticscale",
	"cloudalpha",
	"cloudbox",
	"cloudboxcolour",
	"cloudcolour",
	"cloudfade",
	"cloudheight",
	"cloudlayer",
	"cloudscale",
	"cloudscrollx",
	"cloudscrolly",
	"exec",
	"fog",
	"fogcolour",
	"fogdomecap",
	"fogdomeclip",
	"fogdomeclouds",
	"fogdomecolour",
	"fogdomeheight",
	"fogdomemax",
	"fogdomemin",
	"grassalpha",
	"lightprecision",
	"lmshadows",
	"loadsky",
	"mapmodel",
	"mapmodelreset",
	"mapmsg",
	"mapsound",
	"maptitle",
	"minimapclip",
	"minimapcolour",
	"minimapheight",
	"mmodel",
	"mmodelfile",
	"setshader",
	"setshaderparam",
	"shadowmapambient",
	"shadowmapangle",
	"skybox",
	"skyboxcolour",
	"skylight",
	"skytexture",
	"skytexturelight",
	"spinclouds",
	"spinsky",
	"sunlight",
	"sunlightpitch",
	"sunlightscale",
	"sunlightyaw",
	"texalpha",
	"texcolor",
	"texffenv",
	"texlayer",
	"texoffset",
	"texrotate ",
	"texscale",
	"texscroll",
	"texture",
	"texturereset",
	"water2colour",
	"water2fog",
	"watercolour",
	"waterfallcolour",
	"waterfog",
	"waterspec",
	"yawsky",
}

// Check for any of these characters before further processing a line
var EscapeCharacters = regexp.MustCompile(`(;|\(|\[)`)

// Match leading and trailing quotes in a string
var Quotes = regexp.MustCompile(`^"(.*)"$`)

// A resource is a tuple of the type of the resource (see Property) and it's path
type Resource struct {
	Property, Path string
}

// Reads a map config file and returns all resources declared in the config
func ReadResources(config io.Reader) ([]Resource, error) {
	resources := make([]Resource, 0)
	scanner := bufio.NewScanner(config)

	count := 0
	for scanner.Scan() {
		// Increase line count
		count++
		// Return line from buffer
		line := scanner.Text()
		fields := strings.Split(line, " ")
		property := fields[0]

		if index, ok := Properties[property]; ok {
			if Properties[property] < len(fields) {
				value := Quotes.ReplaceAllString(strings.TrimSpace(fields[index]), `$1`)
				if len(value) > 0 {
					resources = append(resources, Resource{property, value})
				} else {
					return nil, errors.New(fmt.Sprintln("To few arguments in line", count))
				}
			} else {
				return nil, errors.New(fmt.Sprintln("To few arguments in line", count))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return resources, nil
}

// Filters a config file to retain only keys used in map configs
// Bad characters are also stripped away at the potential risk of making commands unusable
func Filter(config io.Reader) (io.Reader, error) {
	output := new(bytes.Buffer)

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		line := scanner.Text()
		// Comments ain't never harm anybody
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			fmt.Fprintln(output, line)
		} else if collection.Any(Whitelist, func (item string) bool {
			return strings.HasPrefix(line, item)
		}) {
			// This prevents the bad boys from conquering the shire
			escaped := EscapeCharacters.Split(line, 2)
			fmt.Fprintln(output, escaped[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return bytes.NewReader(output.Bytes()), nil
}
