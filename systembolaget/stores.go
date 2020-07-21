package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
)

// Stores ...
type Stores struct {
	Info struct {
		Message string `systembolagetXML:"Meddelande" json:"message"`
	} `json:"info"`
	Stores []struct {
		Type         string `systembolagetXML:"Typ" json:"type"`
		ID           string `systembolagetXML:"Nr" json:"id"`
		Name         string `systembolagetXML:"Namn" json:"name"`
		Address1     string `json:address1`
		Address2     string `json:address2`
		Address3     string `json:address3`
		Address4     string `json:address4`
		Address5     string `json:address5`
		PhoneNumber  string `systembolagetXML:"Telefon" json:"phoneNumber"`
		StoreType    string `systembolagetXML:"ButiksTyp" json:"storeType"`
		Services     string `systembolagetXML:"Tjanster" json:"services"`
		Keywords     string `systembolagetXML:"SokOrd" json:"keywords"`
		OpeningHours string `systembolagetXML:"Oppettider" json:"openingHours"`
		RT90x        string `json:"rt90x"`
		RT90y        string `json:"rt90y"`
	} `systembolagetXML:"ButikOmbud" xml:"Store" json:"stores"`
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
	err = util.UnmarshalXML(bytes, &response)
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
