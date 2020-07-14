package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
)

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
type strippedInventory struct {
	Info struct {
		Message string
	}
	Stores []struct {
		ID          string
		ItemNumbers []int
	}
}

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

	return response, nil
}

// ConvertToJSON ...
func (response *Inventory) ConvertToJSON() ([]byte, error) {
	return json.Marshal(response)
}

// ConvertToXML ...
func (response *Inventory) ConvertToXML() ([]byte, error) {
	// Workaround for ignoring the struct tags in the field
	// when marshalling to XML
	strippedResponse := strippedInventory(*response)
	return xml.Marshal(strippedResponse)
}
