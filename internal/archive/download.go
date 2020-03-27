package archive

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"net/url"
	"time"
)

// Downloads a single file and returns the destination path
func DownloadFile(source url.URL, destination string) (string, error) {
	request, err := grab.NewRequest(destination, source.String())
	if err != nil {
		return "", err
	}
	client := grab.NewClient()
	response := client.Do(request)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	Loop:
	for {
		select {
		case <-response.Done:
			// download is complete
			break Loop
		}
	}

	return destination, nil
}

// Downloads a batch of files.
// Returns a slice of destination paths
func DownloadBatch(sources []url.URL, destinations []string) ([]string, error) {
	requests := make([]*grab.Request, 0)
	client := grab.NewClient()

	for index, source := range sources {
		request, err := grab.NewRequest(destinations[index], source.String())
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	responsesChannel := client.DoBatch(5, requests...)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for resp := range responsesChannel {
		/*if err := resp.Err(); err != nil {
			panic(err)
		}*/
		fmt.Printf("Downloaded %s to %s\n", resp.Request.URL(), resp.Filename)
	}

	return destinations, nil
}
