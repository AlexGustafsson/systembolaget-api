package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
	"github.com/jinzhu/copier"
)

type storesInput struct {
	Info struct {
		Message string `xml:"Meddelande"`
	} `json:"info"`
	Stores []struct {
		Type         string `xml:"Typ"`
		ID           string `xml:"Nr"`
		Name         string `xml:"Namn"`
		Address1     string
		Address2     string
		Address3     string
		Address4     string
		Address5     string
		PhoneNumber  string `xml:"Telefon"`
		StoreType    string `xml:"ButiksTyp"`
		Services     string `xml:"Tjanster"`
		Keywords     string `xml:"SokOrd"`
		OpeningHours string `xml:"Oppettider"`
		RT90x        string
		RT90y        string
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
	} `xml:"Store" json:"stores"`
}

// Download ...
func (stores *Stores) Download() error {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/stores/xml")
	if err != nil {
		return err
	}

	// Unmarshal
	var response = &storesInput{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return err
	}

	// Deep copy the response into the similar, but not equal struct
	// This is a workaround for not supporting different struct tags
	// for serialization and deserialization
	copier.Copy(&stores, &response)

	return nil
}

// ParseFromXML ...
func (stores *Stores) ParseFromXML(bytes []byte) error {
	return xml.Unmarshal(bytes, stores)
}

// ParseFromJSON ...
func (stores *Stores) ParseFromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, stores)
}

// ConvertToJSON ...
func (stores *Stores) ConvertToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(stores, "", "  ")
	}

	return json.Marshal(stores)
}

// ConvertToXML ...
func (stores *Stores) ConvertToXML(pretty bool) ([]byte, error) {
	if pretty {
		return xml.MarshalIndent(stores, "", "  ")
	}

	return xml.Marshal(stores)
}
