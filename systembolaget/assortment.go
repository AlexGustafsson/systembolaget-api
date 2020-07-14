package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
)

// Assortment ...
type Assortment struct {
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

// Workaround for ignoring the struct tags in the field
// when marshalling to XML
type strippedAssortment struct {
	Created string
	Info    struct {
		Message string
	}
	Items []struct {
		ID         string
		ItemID     string
		ItemNumber string
		Name       string
		Name2      string
		// Including VAT:
		Price float32
		// In millilitres
		Volume                float32
		PricePerLiter         float32
		SalesStart            string
		Discontinued          bool
		Group                 string
		Type                  string
		Style                 string
		Packaging             string
		Seal                  string
		Origin                string
		CountryOfOrigin       string
		Producer              string
		Supplier              string
		Vintage               string
		TestedVintage         string
		AlcholByVolume        string
		Assortment            string
		AssortmentText        string
		Organic               bool
		Ethical               bool
		Kosher                bool
		IngredientDescription string
	}
}

// DownloadAssortment downloads the data and unmarshals the original XML format.
func DownloadAssortment() (*Assortment, error) {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/products/xml")
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var response = &Assortment{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ConvertToJSON ...
func (response *Assortment) ConvertToJSON() ([]byte, error) {
	return json.Marshal(response)
}

// ConvertToXML ...
func (response *Assortment) ConvertToXML() ([]byte, error) {
	// Workaround for ignoring the struct tags in the field
	// when marshalling to XML
	strippedResponse := strippedAssortment(*response)
	return xml.Marshal(strippedResponse)
}
