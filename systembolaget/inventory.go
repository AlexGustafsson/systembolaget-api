package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
	"github.com/jinzhu/copier"
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
	}
	Stores []struct {
		ID          string `xml:"ButikNr,attr"`
		ItemNumbers []int  `xml:"ArtikelNr"`
	} `xml:"Butik"`
}

// Store ...
type Store struct {
	ID          string `json:"id"`
	ItemNumbers []int  `xml:"ItemNumber" json:"itemNumbers"`
}

// Inventory ...
type Inventory struct {
	Info struct {
		Message string `json:"message"`
	} `json:"info"`
	Stores []struct {
		ID          string `json:"id"`
		ItemNumbers []int  `xml:"ItemNumber" json:"itemNumbers"`
	} `xml:"Store" json:"stores"`
}

type byStoreID []struct {
	ID          string
	ItemNumbers []int
}

func (a byStoreID) Len() int           { return len(a) }
func (a byStoreID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byStoreID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Download ...
func (inventory *Inventory) Download() error {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/stock/xml")
	if err != nil {
		return err
	}

	// Unmarshal
	var response = &InventoryInput{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return err
	}

	// Deep copy the response into the similar, but not equal struct
	// This is a workaround for not supporting different struct tags
	// for serialization and deserialization
	copier.Copy(&inventory, &response)

	// Sort arrays
	sort.Sort(byStoreID(inventory.Stores))
	for _, store := range inventory.Stores {
		sort.Ints(store.ItemNumbers)
	}

	return nil
}

// ParseFromXML ...
func (inventory *Inventory) ParseFromXML(bytes []byte) error {
	return xml.Unmarshal(bytes, inventory)
}

// ParseFromJSON ...
func (inventory *Inventory) ParseFromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, inventory)
}

// ConvertToJSON ...
func (inventory *Inventory) ConvertToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(inventory, "", "  ")
	}

	return json.Marshal(inventory)
}

// ConvertToXML ...
func (inventory *Inventory) ConvertToXML(pretty bool) ([]byte, error) {
	if pretty {
		return xml.MarshalIndent(inventory, "", "  ")
	}

	return xml.Marshal(inventory)
}
