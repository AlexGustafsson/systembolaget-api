package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// AssortmentAPI ...
type AssortmentAPI struct {
	Container AssortmentContainer
	Bytes     []byte
}

// AssortmentContainer ...
type AssortmentContainer struct {
	Created string `xml:"skapad-tid"`
	Info    struct {
		Message string `xml:"meddelande"`
	} `xml:"info"`
	Items []struct {
		ID         string `xml:"nr"`
		ItemID     string `xml:"Artikelid"`
		ItemNumber string `xml:"Varnummer"`
		Name       string `xml:"Namn"`
		Name2      string `xml:"Namn2"`
		// Including VAT:
		Price float32 `xml:"Prisinklmoms"`
		// In millilitres
		Volume                float32 `xml:"Volymiml"`
		PricePerLiter         float32 `xml:"PrisPerLiter"`
		SalesStart            string  `xml:"Saljstart"`
		Discontinued          bool    `xml:"Utgått"`
		Group                 string  `xml:"Varugrupp"`
		Type                  string  `xml:"Typ"`
		Style                 string  `xml:"Stil"`
		Packaging             string  `xml:"Forpackning"`
		Seal                  string  `xml:"Forslutning"`
		Origin                string  `xml:"Ursprung"`
		CountryOfOrigin       string  `xml:"Ursprunglandnamn"`
		Producer              string  `xml:"Producent"`
		Supplier              string  `xml:"Leverantor"`
		Vintage               string  `xml:"Argang"`
		TestedVintage         string  `xml:"Provadargang"`
		AlcholByVolume        string  `xml:"Alkoholhalt"`
		Assortment            string  `xml:"Sortiment"`
		AssortmentText        string  `xml:"SortimentText"`
		Organic               bool    `xml:"Ekologisk"`
		Ethical               bool    `xml:"Etisk"`
		Kosher                bool    `xml:"Koscher"`
		IngredientDescription string  `xml:"RavarorBeskrivning"`
	} `xml:"artikel"`
}

// Debug ...
func (assortmentAPI *AssortmentAPI) Debug() {
	fmt.Printf("There is a message: %s\n", assortmentAPI.Container.Info.Message)

	fmt.Printf("There are %d Assortment\n", len(assortmentAPI.Container.Items))
}

// Download ...
func (assortmentAPI *AssortmentAPI) Download() error {
	bytes, err := download("https://www.systembolaget.se/api/assortment/products/xml")

	assortmentAPI.Bytes = bytes

	xml.Unmarshal(bytes, &assortmentAPI.Container)

	return err
}

// Cache ...
func (assortmentAPI *AssortmentAPI) Cache() error {
	err := ioutil.WriteFile("output/xml/assortment.xml", assortmentAPI.Bytes, 0644)

	return err
}

// Convert ...
func (assortmentAPI *AssortmentAPI) Convert() error {
	bytes, _ := json.Marshal(assortmentAPI.Container)
	err := ioutil.WriteFile("output/json/assortment.json", bytes, 0644)

	return err
}

// Process ...
func (assortmentAPI *AssortmentAPI) Process() error {
	err := assortmentAPI.Download()

	if err != nil {
		return err
	}

	assortmentAPI.Debug()
	err = assortmentAPI.Cache()

	if err != nil {
		return err
	}

	err = assortmentAPI.Convert()

	return err
}
