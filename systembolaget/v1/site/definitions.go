package site

// Site ...
type Site struct {
	OpeningHours string
	IsTastingStore bool
	SiteID string `json:"SiteId"`
	Alias string
	Address string
	DisplayName string
	PostalCode string
	City string
	County string
	Country string
	IsStore bool
	IsAgent bool
	IsActiveForAgentOrder bool
	Phone string
	Email string
	Services string
	Depot string
	Name string
	Position struct{
		Long float32
		Lat float32
	}
}

// Query ...
type Query struct {
  Page *int
  SearchQuery *string
  SiteType *int
}

// SearchResult ...
type SearchResult struct {
  Metadata struct{
    DocCount int
    NextPage int
    Query Query
  }
  Hits []Site
}
