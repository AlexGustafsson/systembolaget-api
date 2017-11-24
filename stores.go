package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// StoresAPI ...
type StoresAPI struct {
	Container StoresContainer
	Bytes     []byte
}

// StoresContainer ...
type StoresContainer struct {
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

// Debug ...
func (storesAPI *StoresAPI) Debug() {
	fmt.Printf("There is a message: %s\n", storesAPI.Container.Info.Message)

	fmt.Printf("There are %d stores\n", len(storesAPI.Container.Stores))
}

// Download ...
func (storesAPI *StoresAPI) Download() error {
	bytes, err := download("https://www.systembolaget.se/api/assortment/stores/xml")

	storesAPI.Bytes = bytes

	xml.Unmarshal(bytes, &storesAPI.Container)

	return err
}

// Cache ...
func (storesAPI *StoresAPI) Cache() error {
	err := ioutil.WriteFile("output/xml/stores.xml", storesAPI.Bytes, 0644)

	return err
}

// Convert ...
func (storesAPI *StoresAPI) Convert() error {
	bytes, _ := json.Marshal(storesAPI.Container)
	err := ioutil.WriteFile("output/json/stores.json", bytes, 0644)

	return err
}

// Process ...
func (storesAPI *StoresAPI) Process() error {
	err := storesAPI.Download()

	if err != nil {
		return err
	}

	storesAPI.Debug()
	err = storesAPI.Cache()

	if err != nil {
		return err
	}

	err = storesAPI.Convert()

	return err
}
