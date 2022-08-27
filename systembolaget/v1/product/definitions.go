package product

// Product ...
type Product struct {
	ProductID                 string `json:"productId"`
	ProductNumber             string
	ProductNameBold           string
	ProductNameThin           string
	Category                  string
	ProductNumberShort        string
	ProducerName              string
	SupplierName              string
	IsKosher                  bool
	BottleTextShort           string
	Seal                      string
	RestrictedParcelQuantity  int
	IsOrganic                 bool
	IsEthical                 bool
	EthicalLabel              string
	IsWebLaunch               bool
	SellStartDate             string
	IsCompletelyOutOfStock    bool
	IsTemporaryOutOfStock     bool
	AlcoholPercentage         float32
	Volume                    float32
	Price                     float32
	Country                   string
	OriginLevel1              string
	OriginLevel2              string
	Vintage                   int
	SubCategory               string
	Type                      string
	Style                     string
	AssortmentText            string
	BeverageDescriptionShort  string
	Usage                     string
	Taste                     string
	Assortment                string
	RecycleFee                float32
	IsManufacturingCountry    bool
	IsRegionalRestricted      bool
	IsInStoreSearchAssortment string
}

// Inventory ...
type Inventory struct {
	SiteID   string `json:"SiteId"`
	Products []struct {
		ProductID     string `json:"ProductId"`
		ProductNumber string
	}
}

// Query ...
type Query struct {
	SortDirection        *int
	SortBy               *int
	Page                 *int
	SearchQuery          *string
	Country              *string
	OriginLevel1         *string
	OriginLevel2         *string
	OriginLevel3         *string
	Vintage              *string
	SubCategory          *string
	Type                 *string
	Style                *string
	AssortmentText       *string
	BottleTypeGroup      *string
	Seal                 *string
	Grape                *string
	PriceMin             *float32
	PriceMax             *float32
	AlcoholPercentageMin *float32
	AlcoholPercentageMax *float32
	News                 *string
	Csr                  *string
	OtherSelections      *string
	SellStartDateFrom    *string
	SellStartDateTo      *string
}

// SearchResult ...
type SearchResult struct {
	Metadata struct {
		DocCount   int
		NextPage   int
		PriceRange struct {
			Min float32
			Max float32
		}
		Query Query
	}
	Hits []Product
}
