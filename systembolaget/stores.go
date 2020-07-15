package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
)

// Stores ...
type Stores struct {
	Info struct {
		Message string `xml:"Meddelande" json:"message"`
	} `json:"info"`
	Stores []struct {
		Type         string `xml:"Typ" json:"type"`
		ID           string `xml:"Nr" json:"id"`
		Name         string `xml:"Namn" json:"name"`
		Address1     string `json:"address1"`
		Address2     string `json:"address2"`
		Address3     string `json:"address3"`
		Address4     string `json:"address4"`
		Address5     string `json:"address5"`
		PhoneNumber  string `xml:"Telefon" json:"phoneNumber"`
		StoreType    string `xml:"ButiksTyp" json:"storeType"`
		Services     string `xml:"Tjanster" json:"services"`
		Keywords     string `xml:"SokOrd" json:"keywords"`
		OpeningHours string `xml:"Oppettider" json:"openingHours"`
		RT90x        string `json:"rt90x"`
		RT90y        string `json:"rt90y"`
	} `xml:"ButikOmbud" json:"stores"`
}

// Workaround for ignoring the struct tags in the field
// when marshalling to XML
type strippedStores struct {
	Info struct {
		Message string
	}
	Stores []struct {
		Type         string
		ID           string
		Name         string
		Address1     string
		Address2     string
		Address3     string
		Address4     string
		Address5     string
		PhoneNumber  string
		StoreType    string
		Services     string
		Keywords     string
		OpeningHours string
		RT90x        string
		RT90y        string
	}
}

// DownloadStores ...
func DownloadStores() (*Stores, error) {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/stores/xml")
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var response = &Stores{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ConvertToJSON ...
func (response *Stores) ConvertToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(response, "", "  ")
	}

	return json.Marshal(response)
}

// ConvertToXML ...
func (response *Stores) ConvertToXML(pretty bool) ([]byte, error) {
	// Workaround for ignoring the struct tags in the field
	// when marshalling to XML
	strippedResponse := strippedStores(*response)
	if pretty {
		return xml.MarshalIndent(strippedResponse, "", "  ")
	}

	return xml.Marshal(strippedResponse)
}
