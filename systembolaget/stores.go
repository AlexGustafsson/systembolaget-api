package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
)

// StoresInput ...
type StoresInput struct {
	Info struct {
		Message string `xml:"Meddelande"`
	} `json:"info"`
	Stores []struct {
		Type         string `xml:"Typ"`
		ID           string `xml:"Nr"`
		Name         string `xml:"Namn"`
		Address1     string `json:"address1"`
		Address2     string `json:"address2"`
		Address3     string `json:"address3"`
		Address4     string `json:"address4"`
		Address5     string `json:"address5"`
		PhoneNumber  string `xml:"Telefon"`
		StoreType    string `xml:"ButiksTyp"`
		Services     string `xml:"Tjanster"`
		Keywords     string `xml:"SokOrd"`
		OpeningHours string `xml:"Oppettider"`
		RT90x        string `json:"rt90x"`
		RT90y        string `json:"rt90y"`
	} `xml:"ButikOmbud"`
}

// Stores ...
type Stores struct {
	Info struct {
		Message string `json:"message"`
	} `json:"info"`
	Stores []struct {
		Type         string `json:"type"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Address1     string `json:"address1"`
		Address2     string `json:"address2"`
		Address3     string `json:"address3"`
		Address4     string `json:"address4"`
		Address5     string `json:"address5"`
		PhoneNumber  string `json:"phoneNumber"`
		StoreType    string `json:"storeType"`
		Services     string `json:"services"`
		Keywords     string `json:"keywords"`
		OpeningHours string `json:"openingHours"`
		RT90x        string `json:"rt90x"`
		RT90y        string `json:"rt90y"`
	} `json:"stores"`
}

// DownloadStores ...
func DownloadStores() (*Stores, error) {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/stores/xml")
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var response = &StoresInput{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var convertedResponse = Stores(*response)
	return &convertedResponse, nil
}

// ParseStoreFromXML ...
func ParseStoreFromXML(bytes []byte) (*Stores, error) {
	var response = &Stores{}
	err := xml.Unmarshal(bytes, &response)
	return response, err
}

// ParseStoreFromJSON ...
func ParseStoreFromJSON(bytes []byte) (*Stores, error) {
	var response = &Stores{}
	err := json.Unmarshal(bytes, &response)
	return response, err
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
	if pretty {
		return xml.MarshalIndent(response, "", "  ")
	}

	return xml.Marshal(response)
}
