package systembolaget

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Range describes a range of numbers.
type Range struct {
	Minimum float32 `json:"min"`
	Maximum float32 `json:"max"`
}

// SearchResult contains the raw result of a query.
type SearchResult struct {
	Metadata struct {
		DocumentCount                int   `json:"docCount"`
		FullAssortmentDocumentCount  int   `json:"fullAssortmentDocCount"`
		NextPage                     int   `json:"nextPage"`
		PriceRange                   Range `json:"priceRange"`
		VolumeRange                  Range `json:"volumeRange"`
		AlcoholPercentageRange       Range `json:"alcoholPercantageRange"`
		SugarContentRange            Range `json:"sugarContentRange"`
		SugarContentGramsPer100Range Range `json:"sugarContentGramPer100mlRange"`
	} `json:"metadata"`
	Products []Product `json:"products"`
	Filters  []Filter  `json:"filters"`
}

// Product describes a product in Systembolaget's assortment.
// The typing is intentionally left vague as it's expected that the types may
// change often, which would lead this API to break.
type Product map[string]any

// Filter is returned in a search result and describes how to further narrow
// the search.
type Filter struct {
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	DisplayName           string `json:"displayName"`
	Description           string `json:"description"`
	Summary               string `json:"symmary"`
	LegalText             string `json:"legalText"`
	IsMultipleChoice      bool   `json:"isMultipleChoice"`
	IsActive              bool   `json:"isActive"`
	IsSubtitleTextVisible bool   `json:"isSubtitleTextVisible"`
	SearchModifiers       []SearchModifier
	Child                 *Filter `json:"child"`
}

// SearchModifier describes modifiers to a filter.
type SearchModifier struct {
	Value        string `json:"value"`
	Count        int    `json:"count"`
	IsActive     bool   `json:"isActive"`
	SubtitleText string `json:"subtitleText"`
	FriendlyURL  string `json:"friendlyUrl"`
}

// SearchCursor allows to easily loop through the results of a query.
type SearchCursor struct {
	client *Client

	options     SearchOptions
	filters     []SearchFilter
	currentPage *SearchResult

	currentError error
	index        int
}

// Next moves the cursor to the next product.
// Returns whether or not there was any more products.
// Use At() and Error() to get the product and potential error.
func (c *SearchCursor) Next(ctx context.Context, delayBetweenPages time.Duration) bool {
	c.index++

	log := GetLogger(ctx).With(slog.Duration("delayBetweenPages", delayBetweenPages))

	// There's at least one product left on the current page
	if c.currentPage != nil && c.index < len(c.currentPage.Products) {
		c.currentError = nil
		return true
	}

	// There are no more pages
	if !c.hasNextPage() {
		log.Debug("No more pages")
		c.index = -1
		c.currentError = nil
		return false
	}

	// Move on to the next page after waiting for the specified amount of time
	// The first page is fetched immediately
	if c.currentPage != nil {
		log.Debug("Waiting")
		<-time.After(delayBetweenPages)
	}
	c.currentError = c.nextPage(ctx, log)
	if c.currentError != nil {
		return false
	}

	return len(c.currentPage.Products) > 0
}

// hasNextPage returns whether or not the cursor has a next page.
func (c *SearchCursor) hasNextPage() bool {
	// currentPage is nil if no page has been fetched
	// nextPage is only set if there is a new page to fetch
	return c.currentPage == nil || c.currentPage.Metadata.NextPage > 0
}

// nextPage fetches the next page.
func (c *SearchCursor) nextPage(ctx context.Context, log *slog.Logger) error {
	c.options.Page++
	c.index = 0

	log = log.With(slog.Int("index", c.index))

	nextPage, err := c.client.Search(SetLogger(ctx, log), &c.options, c.filters...)
	if err != nil {
		return err
	}

	c.currentPage = nextPage
	return nil
}

// At returns the current product. Returns nil if there's no product to retrieve
// which happens if no page has been succesfully fetched or if an error occured.
func (c *SearchCursor) At() Product {
	if c.currentPage == nil || c.index < 0 || c.index >= len(c.currentPage.Products) {
		return nil
	}
	return c.currentPage.Products[c.index]
}

// Error returns the latest error, if any. Populated when pages are fetched.
func (c *SearchCursor) Error() error {
	return c.currentError
}

// SortProperty is a property which can be used to sort search results.
type SortProperty string

const (
	SortPropertyScore             SortProperty = "Score"
	SortPropertyPrice             SortProperty = "Price"
	SortPropertyName              SortProperty = "Name"
	SortPropertyVolume            SortProperty = "Volume"
	SortPropertyProductLaunchDate SortProperty = "ProductLaunchDate"
	SortPropertyVintage           SortProperty = "Vintage"
)

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "Ascending"
	SortDirectionDescending SortDirection = "Descending"
)

// SearchOptions contains optional search options.
type SearchOptions struct {
	// PageSize is the size of pages returned by the server.
	// Note that there's an upper limit of 30 results enforced by the server.
	// Defaults to 30.
	PageSize int
	// Page is the page of results to return. Defaults to 1.
	Page int

	// SortBy specifies the property to sort by, such as "Name".
	SortBy SortProperty
	// SortDirection specifies the direction in which to sort. "Ascending" or "Descending".
	SortDirection SortDirection
}

// SearchFilter describes filters that modifies the search query.
type SearchFilter func(*url.Values)

// FilterByStore filters products that are in a specific store's assortment.
func FilterByStore(store string) SearchFilter {
	return func(v *url.Values) {
		v.Set("storeId", store)
	}
}

// FilterByQuery queries products using free text.
func FilterByQuery(query string) SearchFilter {
	return func(v *url.Values) {
		v.Set("textQuery", query)
	}
}

// FilterByTasteClockBody filters products of a certain body where 0 (minimum) is a thin body and 12 (maximum) is a full body.
func FilterByTasteClockBody(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockBody.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockBody.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockBitterness filters products of a certain bitterness where 0 (minimum) is not bitter at all and 12 (maximum) is very bitter.
func FilterByTasteClockBitterness(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockBitter.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockBitter.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockSweetness filters products of a certain sweetness where 0 (minimum) is not sweet at all and 12 (maximum) is very sweet.
func FilterByTasteClockSweetness(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockSweetness.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockSweetness.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockSmokiness filters products of a certain smokiness where 0 (minimum) is not smoky at all and 12 (maximum) is very smoky.
func FilterByTasteClockSmokiness(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockSmokiness.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockSmokiness.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByVintage may be used more than once.
func FilterByVintage(vintage int) SearchFilter {
	return func(v *url.Values) {
		v.Add("vintage", strconv.FormatInt(int64(vintage), 10))
	}
}

func FilterByProductLaunch(min time.Time, max time.Time) SearchFilter {
	return func(v *url.Values) {
		v.Set("productLaunch.min", min.Format("2006-01-02"))
		v.Set("productLaunch.max", max.Format("2006-01-02"))
	}
}

func FilterByAlcoholPercentage(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("alcoholPercentage.min", strconv.FormatInt(int64(min), 10))
		v.Set("alcoholPercentage.max", strconv.FormatInt(int64(max), 10))
	}
}

func FilterBySugarContent(min float32, max float32) SearchFilter {
	return func(v *url.Values) {
		v.Set("sugarContentGramPer100ml.min", strconv.FormatFloat(float64(min), 'f', 2, 32))
		v.Set("sugarContentGramPer100ml.max", strconv.FormatFloat(float64(max), 'f', 2, 32))
	}
}

// FilterByGrapes may be used more than once.
func FilterByGrapes(grapes string) SearchFilter {
	return func(v *url.Values) {
		v.Add("grapes", grapes)
	}
}

// FilterByMatch filters products that fit with a taste, such as "Aperitif", "Asiatiskt" or "Kött".
// May be used more than once.
func FilterByMatch(match string) SearchFilter {
	return func(v *url.Values) {
		v.Add("tasteSymbols", match)
	}
}

// FilterByAssortment specifies the assortment the product should be included in, such as "Fast sortiment" or "Tillfälligt sortiment".
// May be used more than once.
func FilterByAssortment(assortment string) SearchFilter {
	return func(v *url.Values) {
		v.Add("assortmentText", assortment)
	}
}

// FilterBySeal filters products that use a specific seal, such as "A-koppling" or "Champagnekork-natur".
func FilterBySeal(seal string) SearchFilter {
	return func(v *url.Values) {
		v.Add("seal", seal)
	}
}

func FilterByVolume(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("volume.min", strconv.FormatInt(int64(min), 10))
		v.Set("volume.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByPackaging filters products that use a specific packaging, such as "Flaska" + "Glasflaska".
// Leave subcategory empty to only filter by category.
func FilterByPackaging(category string, subcategory string) SearchFilter {
	return func(v *url.Values) {
		v.Set("packagingLevel1", category)
		if subcategory != "" {
			v.Set("packagingLevel2", subcategory)
		}
	}
}

func FilterByPrice(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("price.min", strconv.FormatInt(int64(min), 10))
		v.Set("price.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByOrigin filters products that originate from a specific country.
// May be used more than once.
func FilterByOrigin(country string) SearchFilter {
	return func(v *url.Values) {
		v.Add("country", country)
	}
}

// FilterByCategory filters products of a specific category, such as "Öl" + "Ljus lager" + "Internationell stil".
// Leave either subcategory empty to only filter by category.
// Can be used more than once to specify multiple subsubcategories of the same category and subcategory.
func FilterByCategory(category string, subcategory string, subsubcategory string) SearchFilter {
	return func(v *url.Values) {
		v.Set("categoryLevel1", category)
		if subcategory != "" {
			v.Set("categoryLevel2", subcategory)
			if subsubcategory != "" {
				v.Add("categoryLevel3", subsubcategory)
			}
		}
	}
}

// Search searches for products.
func (c *Client) Search(ctx context.Context, options *SearchOptions, filters ...SearchFilter) (*SearchResult, error) {
	log := GetLogger(ctx)

	if options.PageSize == 0 {
		options.PageSize = 30
	}
	if options.Page == 0 {
		options.Page = 1
	}

	query := url.Values{}
	query.Set("size", strconv.FormatInt(int64(options.PageSize), 10))
	query.Set("page", strconv.FormatInt(int64(options.Page), 10))

	if options.SortBy != "" {
		query.Set("sortBy", string(options.SortBy))
	}

	if options.SortDirection != "" {
		query.Set("sortDirection", string(options.SortDirection))
	}

	for _, filter := range filters {
		filter(&query)
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "api-extern.systembolaget.se",
		Path:     "/sb-api-ecommerce/v1/productsearch/search",
		RawQuery: query.Encode(),
	}

	log = log.With(slog.String("url", u.String()))

	header := http.Header{}
	header.Set("Origin", "https://www.systembolaget.se")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Pragma", "no-cache")
	header.Set("Accept", "application/json")
	header.Set("Accept-Encoding", "gzip")
	header.Set("Cache-Control", "no-cache")
	header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

	if c.userAgent != "" {
		header.Set("User-Agent", c.userAgent)
	}

	req := (&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: header,
	}).Clone(ctx)

	log.Debug("Performing request")
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("Request failed", slog.Any("error", err))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			log.Error("Failed to create gzip reader", slog.Any("error", err))
			return nil, err
		}
	default:
		reader = res.Body
	}
	defer reader.Close()

	var result SearchResult
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	log.Debug("Got results", slog.Int("results", result.Metadata.DocumentCount), slog.Int("nextPage", result.Metadata.NextPage))
	return &result, nil
}

// SearchWithCursor creates a SearchCursor easily loop over any number of results.
func (c *Client) SearchWithCursor(options *SearchOptions, filters ...SearchFilter) *SearchCursor {
	if options == nil {
		options = &SearchOptions{}
	}
	cursor := &SearchCursor{
		client:  c,
		options: *options,
		filters: filters,
		index:   -1,
	}
	cursor.options.Page = 0
	return cursor
}
