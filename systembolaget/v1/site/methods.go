package site

import (
	"github.com/AlexGustafsson/systembolaget-api/utils"
	"reflect"
  "http"
  "xml"
)

// Fetch ...
func Fetch(key string, id string) (Site, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/site/v1/site/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request")
	}

	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var site = &Site{}
	err = xml.Unmarshal(response, &site)
	if err != nil {
		return err
	}

	return site, nil
}

// FetchAll ...
func FetchAll(key string) ([]Site, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/site/v1/site", nil)
	if err != nil {
		return nil, fmt.Errof("Unable to create request")
	}

	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var site = &[]Site{}
	err = xml.Unmarshal(response, &site)
	if err != nil {
		return err
	}

	return site, nil
}

// Search ...
func Search(key string, query *Query) (SearchResult, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/site/v1/site", nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request")
	}
	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	// Add all set values in the query
	queryType := reflect.ValueOf(query)
	for i := 0; i < queryType.NumField(); i++ {
		// Only set non-null values
		if !queryType.Field(i).isNil() {
			request.Form.Set(queryType.Field(i).Name, *queryType.Field(i).Interface())
		}
	}

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var sites = &[]Site{}
	err = xml.Unmarshal(response, &sites)
	if err != nil {
		return err
	}

	return sites, nil
}
