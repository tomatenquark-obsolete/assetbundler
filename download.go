package main

import (
	"fmt"
	"github.com/tomatenquark/assetbundler/pkg/assetbundler"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	url := argsWithoutProg[0]
	destination, err := assetbundler.DownloadMap(url)
	if err != nil {
		panic(err)
	}
	fmt.Println(destination)
}
