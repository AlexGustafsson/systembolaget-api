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
		Message string `xml:"Meddelande"`
	}
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
