package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
	"sort"
)

// StoreInput ...
type StoreInput struct {
	ID          string `xml:"ButikNr,attr"`
	ItemNumbers []int  `xml:"ArtikelNr"`
}

// InventoryInput ...
type InventoryInput struct {
	Info struct {
		Message string `xml:"Meddelande"`
	} `json:"info"`
	Stores []struct {
		ID          string `xml:"ButikNr,attr"`
		ItemNumbers []int  `xml:"ArtikelNr"`
	} `xml:"Butik"`
}

// Store ...
type Store struct {
	ID          string `json:"id"`
	ItemNumbers []int  `json:"itemNumbers"`
}

// Inventory ...
type Inventory struct {
	Info struct {
		Message string `json:"message"`
	} `json:"info"`
	Stores []struct {
		ID          string `json:"id"`
		ItemNumbers []int  `json:"itemNumbers"`
	} `json:"stores"`
}

type byID []struct {
	ID          string
	ItemNumbers []int
}

func (a byID) Len() int           { return len(a) }
func (a byID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// DownloadInventory ...
func DownloadInventory() (*Inventory, error) {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/stock/xml")
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var response = &InventoryInput{}
	err = xml.Unmarshal(bytes, response)
	if err != nil {
		return nil, err
	}

	// Sort arrays
	sort.Sort(byID(response.Stores))
	for _, store := range response.Stores {
		sort.Ints(store.ItemNumbers)
	}

	var convertedResponse = Inventory(*response)
	return &convertedResponse, nil
}

// ParseInventoryFromXML ...
func ParseInventoryFromXML(bytes []byte) (*Inventory, error) {
	var response = &Inventory{}
	err := xml.Unmarshal(bytes, &response)
	return response, err
}

// ParseInventoryFromJSON ...
func ParseInventoryFromJSON(bytes []byte) (*Inventory, error) {
	var response = &Inventory{}
	err := json.Unmarshal(bytes, &response)
	return response, err
}

// ConvertToJSON ...
func (response *Inventory) ConvertToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(response, "", "  ")
	}

	return json.Marshal(response)
}

// ConvertToXML ...
func (response *Inventory) ConvertToXML(pretty bool) ([]byte, error) {
	if pretty {
		return xml.MarshalIndent(response, "", "  ")
	}

	return xml.Marshal(response)
}
