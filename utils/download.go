package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Download performs a HTTP request.
// It handles logging and verification of status codes etc.
func Download(request *http.Request) ([]byte, error) {
	log.Debugf("Downloading from URL %s", request.URL)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Unable to receive response: %v", err)
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

	log.Debugf("Successfully downloaded data")
	return data, nil
}
