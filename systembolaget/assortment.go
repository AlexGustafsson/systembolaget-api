package systembolaget

import (
	"encoding/json"
	"encoding/xml"
	"github.com/alexgustafsson/systembolaget-api/utils"
	"github.com/jinzhu/copier"
	"sort"
)

// InputItem ...
type InputItem struct {
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
}

// AssortmentInput ...
type AssortmentInput struct {
	Created string `xml:"skapad-tid"`
	Info    struct {
		Message string `xml:"meddelande"`
	} `xml:"info"`
	Items []InputItem `xml:"artikel"`
}

// Item ...
type Item struct {
	ID         string `json:"id"`
	ItemID     string `json:"itemId"`
	ItemNumber string `json:"itemNumber"`
	Name       string `json:"name"`
	Name2      string `json:"name2"`
	// Including VAT:
	Price float32 `json:"price"`
	// In millilitres
	Volume                float32 `json:"volume"`
	PricePerLiter         float32 `json:"pricePerLiter"`
	SalesStart            string  `json:"salesStart"`
	Discontinued          bool    `json:"discontinued"`
	Group                 string  `json:"group"`
	Type                  string  `json:"type"`
	Style                 string  `json:"style"`
	Packaging             string  `json:"packaging"`
	Seal                  string  `json:"seal"`
	Origin                string  `json:"origin"`
	CountryOfOrigin       string  `json:"countryOfOrigin"`
	Producer              string  `json:"producer"`
	Supplier              string  `json:"supplier"`
	Vintage               string  `json:"vintage"`
	TestedVintage         string  `json:"testedVintage"`
	AlcholByVolume        string  `json:"alcoholByVolume"`
	Assortment            string  `json:"assortment"`
	AssortmentText        string  `json:"assortmentText"`
	Organic               bool    `json:"organic"`
	Ethical               bool    `json:"ethical"`
	Kosher                bool    `json:"kosher"`
	IngredientDescription string  `json:"ingredientDescription"`
}

// Assortment ...
type Assortment struct {
	Created string `json:"created"`
	Info    struct {
		Message string `json:"message"`
	} `json:"info"`
	Items []Item `json:"items"`
}

type byItemID []Item

func (a byItemID) Len() int           { return len(a) }
func (a byItemID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byItemID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Download ...
func (assortment *Assortment) Download() error {
	// Download
	bytes, err := utils.Download("https://www.systembolaget.se/api/assortment/products/xml")
	if err != nil {
		return err
	}

	// Unmarshal
	var response = &AssortmentInput{}
	err = xml.Unmarshal(bytes, &response)
	if err != nil {
		return err
	}

	// Deep copy the response into the similar, but not equal struct
	// This is a workaround for not supporting different struct tags
	// for serialization and deserialization
	copier.Copy(&assortment, &response)

	// Sort arrays
	sort.Sort(byItemID(assortment.Items))

	return nil
}

// ParseFromXML ...
func (assortment *Assortment) ParseFromXML(bytes []byte) error {
	return xml.Unmarshal(bytes, assortment)
}

// ParseFromJSON ...
func (assortment *Assortment) ParseFromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, assortment)
}

// ConvertToJSON ...
func (assortment *Assortment) ConvertToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(assortment, "", "  ")
	}

	return json.Marshal(assortment)
}

// ConvertToXML ...
func (assortment *Assortment) ConvertToXML(pretty bool) ([]byte, error) {
	if pretty {
		return xml.MarshalIndent(assortment, "", "  ")
	}

	return xml.Marshal(assortment)
}
