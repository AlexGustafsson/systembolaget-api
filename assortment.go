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
	Created string `xml:"skapad-tid" json:"created"`
	Info    struct {
		Message string `xml:"meddelande" json:"message"`
	} `xml:"info" json:"info"`
	Items []struct {
		ID         string `xml:"nr" json:"id"`
		ItemID     string `xml:"Artikelid" json:"itemId"`
		ItemNumber string `xml:"Varnummer" json:"itemNumber"`
		Name       string `xml:"Namn" json:"name"`
		Name2      string `xml:"Namn2" json:"name2"`
		// Including VAT:
		Price float32 `xml:"Prisinklmoms" json:"price"`
		// In millilitres
		Volume                float32 `xml:"Volymiml" json:"volume"`
		PricePerLiter         float32 `xml:"PrisPerLiter" json:"pricePerLiter"`
		SalesStart            string  `xml:"Saljstart" json:"salesStart"`
		Discontinued          bool    `xml:"Utgått" json:"discontinued"`
		Group                 string  `xml:"Varugrupp" json:"group"`
		Type                  string  `xml:"Typ" json:"type"`
		Style                 string  `xml:"Stil" json:"style"`
		Packaging             string  `xml:"Forpackning" json:"packaging"`
		Seal                  string  `xml:"Forslutning" json:"seal"`
		Origin                string  `xml:"Ursprung" json:"origin"`
		CountryOfOrigin       string  `xml:"Ursprunglandnamn" json:"countryOfOrigin"`
		Producer              string  `xml:"Producent" json:"producer"`
		Supplier              string  `xml:"Leverantor" json:"supplier"`
		Vintage               string  `xml:"Argang" json:"vintage"`
		TestedVintage         string  `xml:"Provadargang" json:"testedVintage"`
		AlcholByVolume        string  `xml:"Alkoholhalt" json:"alcoholByVolume"`
		Assortment            string  `xml:"Sortiment" json:"assortment"`
		AssortmentText        string  `xml:"SortimentText" json:"assortmentText"`
		Organic               bool    `xml:"Ekologisk" json:"organic"`
		Ethical               bool    `xml:"Etisk" json:"ethical"`
		Kosher                bool    `xml:"Koscher" json:"kosher"`
		IngredientDescription string  `xml:"RavarorBeskrivning" json:"ingridentDescription"`
	} `xml:"artikel" json:"items"`
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
