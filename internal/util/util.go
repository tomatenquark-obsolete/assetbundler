package util

import (
	"os"
	"path"
)

func GetProjectDirectory() string {
	goPath := os.Getenv("GOPATH")
	return path.Join(goPath, "src/github.com/tomatenquark/assetbundler")
}
