package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
	"sort"
)

// Store ...
type Store struct {
	ID          string `xml:"ButikNr,attr" json:"id"`
	ItemNumbers []int  `xml:"ArtikelNr" json:"itemNumbers"`
}

// Inventory ...
type Inventory struct {
	Info struct {
		Message string `xml:"Meddelande" json:"message"`
	} `json:"info"`
	Stores []struct {
		ID          string `xml:"ButikNr,attr" json:"id"`
		ItemNumbers []int  `xml:"ArtikelNr" json:"itemNumbers"`
	} `xml:"Butik" json:"stores"`
}

// Workaround for ignoring the struct tags in the field
// when marshalling to XML
type strippedStore struct {
	ID          string `xml:"ButikNr,attr" json:"id"`
	ItemNumbers []int  `xml:"ArtikelNr" json:"itemNumbers"`
}

type strippedInventory struct {
	Info struct {
		Message string
	}
	Stores []struct {
		ID          string
		ItemNumbers []int
	}
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
	var response = &Inventory{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	// Sort arrays
	sort.Sort(byID(response.Stores))
	for _, store := range response.Stores {
		sort.Ints(store.ItemNumbers)
	}

	return response, nil
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
	// Workaround for ignoring the struct tags in the field
	// when marshalling to XML
	strippedResponse := strippedInventory(*response)
	if pretty {
		return xml.MarshalIndent(strippedResponse, "", "  ")
	}

	return xml.Marshal(strippedResponse)
}
