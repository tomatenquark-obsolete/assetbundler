package main

import (
	"fmt"
	"github.com/tomatenquark/assetbundler/pkg/assetbundler"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	url := argsWithoutProg[0]
	destination := assetbundler.DownloadMap(url)
	fmt.Println(destination)
}
