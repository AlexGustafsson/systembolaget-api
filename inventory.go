package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

/*
<ButikArtikel xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<Info><Meddelande>Läs användarvillkoren på www.systembolaget.se</Meddelande></Info>
<Butik ButikNr="0102">
<ArtikelNr>176103</ArtikelNr>
*/

// InventoryAPI ...
type InventoryAPI struct {
	Container InventoryContainer
	Bytes     []byte
}

// InventoryContainer ...
type InventoryContainer struct {
	Info struct {
		Message string `xml:"Meddelande"`
	}
	Stores []struct {
		ID          string `xml:"ButikNr,attr"`
		ItemNumbers []int  `xml:"ArtikelNr"`
	} `xml:"Butik"`
}

// Debug ...
func (inventoryAPI *InventoryAPI) Debug() {
	fmt.Printf("There is a message: %s\n", inventoryAPI.Container.Info.Message)

	fmt.Printf("There are %d stores with inventory\n", len(inventoryAPI.Container.Stores))
}

// Download ...
func (inventoryAPI *InventoryAPI) Download() error {
	bytes, err := download("https://www.systembolaget.se/api/assortment/stock/xml")

	inventoryAPI.Bytes = bytes

	xml.Unmarshal(bytes, &inventoryAPI.Container)

	return err
}

// Cache ...
func (inventoryAPI *InventoryAPI) Cache() error {
	err := ioutil.WriteFile("output/xml/inventory.xml", inventoryAPI.Bytes, 0644)

	return err
}

// Convert ...
func (inventoryAPI *InventoryAPI) Convert() error {
	bytes, _ := json.Marshal(inventoryAPI.Container)
	err := ioutil.WriteFile("output/json/inventory.json", bytes, 0644)

	return err
}

// Process ...
func (inventoryAPI *InventoryAPI) Process() error {
	err := inventoryAPI.Download()

	if err != nil {
		return err
	}

	inventoryAPI.Debug()
	err = inventoryAPI.Cache()

	if err != nil {
		return err
	}

	err = inventoryAPI.Convert()

	return err
}
