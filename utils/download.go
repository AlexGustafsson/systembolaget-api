package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Download downloads a URL and returns the received bytes.
// It handles logging of status codes etc.
func Download(url string) ([]byte, error) {
	log.Debugf("Downloading from URL %s", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request: %v", err)
	}

	defer response.Body.Close()

	log.Debugf("Got status code %d", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to complete request. Status code: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading body: %v", err)
	}

	log.Debugf("Successfully downloaded data for URL %s", url)
	return data, nil
}
