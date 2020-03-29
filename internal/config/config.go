package config

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// Check for any of these characters before further processing a line
var EscapeCharacters = regexp.MustCompile(`(;|\(|\[)`)

// Match leading and trailing quotes in a string
var Quotes = regexp.MustCompile(`^"([^"]*)"$`)

// A resource is a tuple of the type of the resource (see Property) and it's path
type FileToLoad struct {
	Command, Path string
}

// Reads a map config file and returns all resources declared in the config
func FilesToLoadFromCfg(config io.Reader) ([]FileToLoad, error) {
	files := make([]FileToLoad, 0)
	scanner := bufio.NewScanner(config)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			// skip empty lines
			continue
		}

		fields := strings.Fields(line)
		cmd := fields[0]

		hasFilePath, filePathPos := HasFilePathParam(cmd)
		if !hasFilePath {
			// skip commands that don't have a file path parameter
			continue
		}

		if len(fields) <= filePathPos {
			return nil, tooFewArgs(cmd, lineNumber)
		}

		filePath := strings.TrimSpace(fields[filePathPos])
		if strings.HasPrefix(filePath, `"`) && strings.HasSuffix(filePath, `"`) {
			filePath = filePath[1 : len(filePath)-1]
		}
		if len(filePath) == 0 {
			return nil, emptyFilePath(cmd, lineNumber)
		}

		files = append(files, FileToLoad{cmd, filePath})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

// Filters a config file to retain only keys used in map configs
// Bad characters are also stripped away at the potential risk of making commands unusable
func Filter(config string) (string, error) {
	var output strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(config))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			// keep comments and empty lines as-is
			output.WriteString(line)
			output.WriteRune('\n')
			continue
		}

		cmd := strings.Fields(line)[0]
		if !ValidInMapCfg(cmd) {
			// skip invalid commands
			continue
		}

		output.WriteString(line)
		output.WriteRune('\n')
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return output.String(), nil
}
