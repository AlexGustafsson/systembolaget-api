package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v2"
)

var appHelpTemplate = `Usage: {{.Name}} [options] [arguments]
{{.Usage}}
Options:
	{{range .Flags}}{{.}}
	{{end}}
Commands:
	{{range .Commands}}{{.Name}}{{ "\t" }}{{.Usage}}
	{{end}}
`

var commandHelpTemplate = `Usage: systembolaget {{.Name}} [options] {{if .ArgsUsage}}{{.ArgsUsage}}{{end}}
{{.Usage}}{{if .Description}}
Description:
	{{.Description}}{{end}}{{if .Flags}}
Options:{{range .Flags}}
	{{.}}{{end}}{{end}}
`

func main() {
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "Interact with Systembolagets APIs"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logs",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "assortment",
			Usage:  "Fetch products",
			Action: ActionAssortment,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "api-key",
					Aliases: []string{"k"},
					Usage:   "API key to use. Defaults to automatically fetching one",
				},
				&cli.PathFlag{
					Name:      "output",
					Aliases:   []string{"o"},
					Usage:     "Path to output",
					TakesFile: true,
				},
				&cli.DurationFlag{
					Name:  "page-delay",
					Usage: "Delay between pages. Defaults to 0s",
					Value: 0,
				},
				&cli.IntFlag{
					Name:  "limit",
					Usage: "Number of results to fetch. Use 0 to return all",
					Value: 0,
				},
				&EnumFlag{
					Name:  "sort-by",
					Usage: "Property to sort by. Defaults to Name",
					Value: EnumValue{
						Value: string(systembolaget.SortPropertyName),
						Choices: []string{
							"Score",
							"Price",
							"Name",
							"Volume",
							"ProductLaunchDate",
							"Vintage",
						},
					},
				},
				&EnumFlag{
					Name:  "sort",
					Usage: "Sort direction. Defaults to ascending",
					Value: EnumValue{
						Value: "ascending",
						Choices: []string{
							"ascending",
							"descending",
						},
					},
				},
				// FilterByStore
				&cli.StringFlag{
					Name:  "store",
					Usage: "Filter products that are in a specific store's assortment.",
				},
				// FilterByQuery
				&cli.StringFlag{
					Name:  "query",
					Usage: "Query products using free text.",
				},
				// FilterByTasteClockBody
				&RangeFlag[int]{
					Name:  "taste-clock-body",
					Usage: "Filters products of a certain body where 0 (minimum) is a thin body and 12 (maximum) is a full body.",
				},
				// FilterByTasteClockBitterness
				&RangeFlag[int]{
					Name:  "taste-clock-bitterness",
					Usage: "Filter products of a certain bitterness where 0 (minimum) is not bitter at all and 12 (maximum) is very bitter.",
				},
				// FilterByTasteClockSweetness
				&RangeFlag[int]{
					Name:  "taste-clock-sweetness",
					Usage: "Filter products of a certain sweetness where 0 (minimum) is not sweet at all and 12 (maximum) is very sweet.",
				},
				// FilterByTasteClockSmokiness
				&RangeFlag[int]{
					Name:  "taste-clock-smokiness",
					Usage: "Filter products of a certain smokiness where 0 (minimum) is not smoky at all and 12 (maximum) is very smoky.",
				},
				// FilterByVintage
				&cli.IntSliceFlag{
					Name: "vintage",
				},
				// FilterByProductLaunch
				&RangeFlag[string]{
					Name: "product-launch",
				},
				// FilterByAlcoholPercentage
				&RangeFlag[int]{
					Name: "alcohol-percentage",
				},
				// FilterBySugarContent
				&RangeFlag[int]{
					Name: "sugar-content",
				},
				// FilterByGrapes
				&cli.StringSliceFlag{
					Name: "grapes",
				},
				// FilterByMatch
				&cli.StringSliceFlag{
					Name:  "match",
					Usage: "Filter products that fit with a taste, such as 'Aperitif', 'Asiatiskt' or 'Kött'.",
				},
				// FilterByAssortment
				&cli.StringSliceFlag{
					Name:  "assortment",
					Usage: "Assortment the product should be included in, such as 'Fast sortiment' or 'Tillfälligt sortiment'",
				},
				// FilterBySeal
				&cli.StringSliceFlag{
					Name:  "seal",
					Usage: "Seal the product should use, such as 'A-koppling' or 'Champagnekork-natur'",
				},
				// FilterByVolume
				&RangeFlag[int]{
					Name: "volume",
				},
				// FilterByPackaging
				&cli.StringFlag{
					Name:  "packaging-category",
					Usage: "Filter products that use a specific packaging, such as 'Flaska'.",
				},
				&cli.StringFlag{
					Name:  "packaging-subcategory",
					Usage: "Filter products that use a specific packaging, such as 'Glasflaska'.",
				},
				// FilterByPrice
				&RangeFlag[int]{
					Name: "price",
				},
				// FilterByOrigin
				&cli.StringSliceFlag{
					Name: "origin",
				},
				// FilterByCategory
				&cli.StringFlag{
					Name: "category",
				},
				&cli.StringFlag{
					Name: "subcategory",
				},
				&cli.StringSliceFlag{
					Name: "subsubcategory",
				},
			},
		},
		{
			Name:   "stores",
			Usage:  "Fetch stores",
			Action: ActionStores,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "api-key",
					Aliases: []string{"k"},
					Usage:   "API key to use. Defaults to automatically fetching one",
				},
				&cli.PathFlag{
					Name:      "output",
					Aliases:   []string{"o"},
					Usage:     "Path to output",
					TakesFile: true,
				},
				&cli.StringFlag{
					Name:    "search",
					Aliases: []string{"q"},
					Usage:   "Optional search query",
				},
			},
		},
		{
			Name:   "stock",
			Usage:  "Fetch stock balance for a specific product in a store",
			Action: ActionStock,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "api-key",
					Aliases: []string{"k"},
					Usage:   "API key to use. Defaults to automatically fetching one",
				},
				&cli.StringFlag{
					Name:     "store-id",
					Usage:    "The ID of the store",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "product-number",
					Usage:    "The product number",
					Required: true,
				},
				&cli.PathFlag{
					Name:      "output",
					Aliases:   []string{"o"},
					Usage:     "Path to output",
					TakesFile: true,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
