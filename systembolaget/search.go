package systembolaget

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
		DocumentCount               int `json:"docCount"`
		FullAssortmentDocumentCount int `json:"fullAssortmentDocCount"`
		// PreviousPage is -1 if there are no more pages.
		PreviousPage int `json:"previousPage"`
		// NextPage is -1 if there are no more pages. It can be a lower value than
		// the requested page if the requested page is out of bounds.
		NextPage                     int    `json:"nextPage"`
		TotalPages                   int    `json:"totalPages"`
		PriceRange                   Range  `json:"priceRange"`
		VolumeRange                  Range  `json:"volumeRange"`
		AlcoholPercentageRange       Range  `json:"alcoholPercantageRange"`
		SugarContentRange            Range  `json:"sugarContentRange"`
		SugarContentGramsPer100Range Range  `json:"sugarContentGramPer100mlRange"`
		DidYouMeanQuery              string `json:"didYouMeanQuery"`
	} `json:"metadata"`
	Products []Product `json:"products"`
	Filters  []Filter  `json:"filters"`
}

// Product describes a product in Systembolaget's assortment.
// The typing is intentionally left vague as it's expected that the types may
// change often, which would lead this API to break.
type Product map[string]any

func (p Product) getNonEmptyString(key string) (string, bool) {
	v, ok := p[key].(string)
	if !ok || v == "" {
		return "", false
	}

	return v, true
}

func (p Product) getNonEmptyFloat64(key string) (float64, bool) {
	v, ok := p[key].(float64)
	if !ok || v == 0 {
		return 0, false
	}

	return v, true
}

// ID returns the product's product id.
func (p Product) ID() (string, bool) {
	return p.getNonEmptyString("productId")
}

// Number returns the product's article number, typically used in stores for
// quick lookup.
func (p Product) Number() (string, bool) {
	return p.getNonEmptyString("productNumber")
}

// Category returns a text describing the product's category.
func (p Product) Category() (string, bool) {
	return p.getNonEmptyString("customCategoryTitle")
}

// Title returns a text describing the product's main name (product name bold).
func (p Product) Title() (string, bool) {
	return p.getNonEmptyString("productNameBold")
}

// Subtitle returns a text describing the product's secondary name (product name
// thin).
func (p Product) Subtitle() (string, bool) {
	return p.getNonEmptyString("productNameThin")
}

// Usage returns a text describing usage of the product (e.g. serving
// temperature and food pairings).
func (p Product) Usage() (string, bool) {
	return p.getNonEmptyString("usage")
}

// Country of origin.
func (p Product) Country() (string, bool) {
	return p.getNonEmptyString("country")
}

// Volume (milliliters).
func (p Product) Volume() (float64, bool) {
	return p.getNonEmptyFloat64("volume")
}

// VolumeText returns a human-readable text describing the product's volume.
func (p Product) VolumeText() (string, bool) {
	return p.getNonEmptyString("volumeText")
}

// Price (SEK).
func (p Product) Price() (float64, bool) {
	return p.getNonEmptyFloat64("price")
}

func (p Product) AlcoholPercentage() (float64, bool) {
	return p.getNonEmptyFloat64("alcoholPercentage")
}

// Thumbnail returns a byte slice typically containing a WebP image.
func (p Product) Thumbnail() ([]byte, bool) {
	imageModules, ok := p["imageModules"].(map[string]any)
	if !ok || imageModules == nil {
		return nil, false
	}

	thumbnail, ok := imageModules["thumbnail"].(string)
	if !ok || thumbnail == "" {
		return nil, false
	}

	data, err := base64.StdEncoding.DecodeString(thumbnail)
	if err != nil {
		return nil, false
	}

	return data, true
}

type ProductImage struct {
	URL string
	// Ignored fields:
	// FileType string `json:"fileType"` // Typically png, but rather unimportant?
	// Size ?? `json:"size"` // Always null?
}

// Images returns references to images of the product.
func (p Product) Images() ([]ProductImage, bool) {
	images, ok := p["images"].([]any)
	if !ok || images == nil {
		return nil, false
	}

	result := make([]ProductImage, 0)
	for _, image := range images {
		i, ok := image.(map[string]any)
		if !ok {
			return nil, false
		}

		imageUrl, ok := i["imageUrl"].(string)
		if !ok || imageUrl == "" {
			continue
		}

		result = append(result, ProductImage{URL: imageUrl})
	}

	return result, true
}

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
	IsExpanded            bool   `json:"isExpanded"`
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
	client *AuthenticatedClient

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

	// There's at least one product left on the current page
	if c.currentPage != nil && c.index < len(c.currentPage.Products) {
		c.currentError = nil
		return true
	}

	// There are no more pages
	if !c.hasNextPage() {
		c.index = -1
		c.currentError = nil
		return false
	}

	// Move on to the next page after waiting for the specified amount of time
	// The first page is fetched immediately
	if c.currentPage != nil {
		select {
		case <-time.After(delayBetweenPages):
		case <-ctx.Done():
			c.currentError = ctx.Err()
			return false
		}
	}
	c.currentError = c.nextPage(ctx)
	if c.currentError != nil {
		return false
	}

	return len(c.currentPage.Products) > 0
}

// hasNextPage returns whether or not the cursor has a next page.
func (c *SearchCursor) hasNextPage() bool {
	// currentPage is nil if no page has been fetched
	// nextPage in the metadata is -1 if there are no more pages
	return c.currentPage == nil || c.currentPage.Metadata.NextPage > 0
}

// nextPage fetches the next page.
func (c *SearchCursor) nextPage(ctx context.Context) error {
	c.options.Page++
	c.index = 0

	nextPage, err := c.client.Search(ctx, &c.options, c.filters...)
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

// CurrentPage returns the current page (if any).
func (c *SearchCursor) CurrentPage() *SearchResult {
	return c.currentPage
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
	// NOTE: The server enforces an upper limit of 30.
	// Defaults to 30.
	PageSize int
	// Page is the page of results to return. Defaults to 1.
	// NOTE: Starts at 1.
	// NOTE: The API seems to limit the number of pages to 333.
	// SEE: https://github.com/AlexGustafsson/systembolaget-api/issues/9
	Page int

	// SortBy specifies the property to sort by, such as "Name".
	SortBy SortProperty
	// SortDirection specifies the direction in which to sort. "Ascending" or
	// "Descending".
	SortDirection SortDirection
}

// SearchFilter describes filters that modifies the search query.
type SearchFilter func(*url.Values)

// FilterByStore filters products that are in a specific store's assortment.
func FilterByStore(store string) SearchFilter {
	return func(v *url.Values) {
		v.Set("storeId", store)
		v.Set("isInStoreAssortmentSearch", "true")
	}
}

// FilterByQuery queries products using free text.
func FilterByQuery(query string) SearchFilter {
	return func(v *url.Values) {
		v.Set("textQuery", query)
	}
}

// FilterByTasteClockBody filters products of a certain body where 0 (minimum)
// is a thin body and 12 (maximum) is a full body.
func FilterByTasteClockBody(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockBody.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockBody.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockBitterness filters products of a certain bitterness where 0
// (minimum) is not bitter at all and 12 (maximum) is very bitter.
func FilterByTasteClockBitterness(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockBitter.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockBitter.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockSweetness filters products of a certain sweetness where 0
// (minimum) is not sweet at all and 12 (maximum) is very sweet.
func FilterByTasteClockSweetness(min int, max int) SearchFilter {
	return func(v *url.Values) {
		v.Set("tasteClockSweetness.min", strconv.FormatInt(int64(min), 10))
		v.Set("tasteClockSweetness.max", strconv.FormatInt(int64(max), 10))
	}
}

// FilterByTasteClockSmokiness filters products of a certain smokiness where 0
// (minimum) is not smoky at all and 12 (maximum) is very smoky.
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

// FilterByMatch filters products that fit with a taste, such as "Aperitif",
// "Asiatiskt" or "Kött".
// May be used more than once.
func FilterByMatch(match string) SearchFilter {
	return func(v *url.Values) {
		v.Add("tasteSymbols", match)
	}
}

// FilterByAssortment specifies the assortment the product should be included
// in, such as "Fast sortiment" or "Tillfälligt sortiment".
// May be used more than once.
func FilterByAssortment(assortment string) SearchFilter {
	return func(v *url.Values) {
		v.Add("assortmentText", assortment)
	}
}

// FilterBySeal filters products that use a specific seal, such as "A-koppling"
// or "Champagnekork-natur".
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

// FilterByPackaging filters products that use a specific packaging, such as
// "Flaska" + "Glasflaska".
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

// FilterByCategory filters products of a specific category, such as "Öl" +
// "Ljus lager" + "Internationell stil".
// Leave either subcategory empty to only filter by category.
// Can be used more than once to specify multiple subsubcategories of the same
// category and subcategory.
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
func (c *AuthenticatedClient) Search(ctx context.Context, options *SearchOptions, filters ...SearchFilter) (*SearchResult, error) {
	if options == nil {
		options = &SearchOptions{}
	}

	if options.PageSize == 0 {
		options.PageSize = 30
	}
	if options.Page < 1 {
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

	header := make(http.Header)
	header.Set("Origin", "https://www.systembolaget.se")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Pragma", "no-cache")
	header.Set("Accept", "application/json")
	header.Set("Cache-Control", "no-cache")
	header.Set("Ocp-Apim-Subscription-Key", c.APIKey)

	if c.UserAgent != "" {
		header.Set("User-Agent", c.UserAgent)
	}

	req := (&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: header,
	}).Clone(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	var result SearchResult
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SearchWithCursor creates a SearchCursor to easily loop over any number of
// results.
func (c *AuthenticatedClient) SearchWithCursor(options *SearchOptions, filters ...SearchFilter) *SearchCursor {
	if options == nil {
		options = &SearchOptions{}
	}
	cursor := &SearchCursor{
		client:  c,
		options: *options,
		filters: filters,
		index:   -1,
	}
	return cursor
}
