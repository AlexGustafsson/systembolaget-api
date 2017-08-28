package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// API ...
type API interface {
	Download() error
	Cache() error
	Convert() error
	Debug()
	Process() error
}

func download(url string) ([]byte, error) {
	fmt.Printf("Downloading %s\n", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Download body: %v", err)
	}

	defer response.Body.Close()

	fmt.Printf("Got status code %d\n", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	fmt.Printf("Successfully read download data for %s \n", url)
	return data, nil
}
