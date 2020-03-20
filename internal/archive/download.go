package archive

import (
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
	client.Do(request)
	return destination, nil
}

// Downloads a batch of files and reports to the progress channel the number of completed requests
// Returns a slice of destination paths
func DownloadBatch(sources []url.URL, destinations []string, progress chan int) ([]string, error) {
	requests := make([]*grab.Request, 0)
	client := grab.NewClient()

	for index, source := range sources {
		request, err := grab.NewRequest(destinations[index], source.String())
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	responses := make([]*grab.Response, 0)
	responsesChannel := client.DoBatch(5, requests...)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for _ = range ticker.C {
		completed := 0
		response, done := <- responsesChannel
		if done {
			completed = len(responses)
			progress <- completed
			break
		} else {
			responses = append(responses, response)
			for _, response := range responses {
				if response.IsComplete() {
					completed++
				}
			}
		}
		progress <- completed
	}

	return destinations, nil
}
