package product

import (
	"github.com/AlexGustafsson/systembolaget-api/utils"
	"reflect"
  "net/http"
  "encoding/xml"
  "fmt"
)

// Fetch ...
func Fetch(key string, id string) (*Product, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/product/v1/product/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request")
	}

	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var product = &Product{}
	err = xml.Unmarshal(response, &product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// FetchAll ...
func FetchAll(key string) ([]Product, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/product/v1/product", nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request")
	}

	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var product = &[]Product{}
	err = xml.Unmarshal(response, &product)
	if err != nil {
		return err
	}

	return product, nil
}

// FetchInventory ...
func FetchInventory(key string) ([]Inventory, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/product/v1/product/getproductswithstore", nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request")
	}

	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	response, err := utils.Download(request)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var inventory = &Inventory{}
	err = xml.Unmarshal(response, &inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

// Search ...
func Search(key string, query *Query) (*SearchResult, error) {
	request, err := http.NewRequest("GET", "https://api-extern.systembolaget.se/product/v1/product/search", nil)
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
	var products = &[]Product{}
	err = xml.Unmarshal(response, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
