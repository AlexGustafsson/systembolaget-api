package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v3"
)

var rootCommandHelpTemplate = `Usage: {{.Name}} [options] [arguments]
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
	cli.RootCommandHelpTemplate = rootCommandHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	root := cli.Command{
		Name:        filepath.Base(os.Args[0]),
		Usage:       "Interact with Systembolagets APIs",
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logs",
			},
		},
		Commands: []*cli.Command{
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
					&cli.DurationFlag{
						Name:  "page-delay",
						Usage: "Delay between pages",
						Value: 0,
					},
					&cli.IntFlag{
						Name:        "limit",
						Usage:       "Number of results to fetch",
						Value:       0,
						DefaultText: "return all",
					},
					&EnumFlag{
						Name:  "sort-by",
						Usage: "Property to sort by",
						Value: string(systembolaget.SortPropertyName),
						Config: EnumConfig{
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
						Value: "ascending",
						Config: EnumConfig{
							Choices: []string{
								"ascending",
								"descending",
							},
						},
					},
					// FilterByStore
					&cli.StringFlag{
						Name:  "store",
						Usage: "Filter products that are in a specific store's assortment",
					},
					// FilterByQuery
					&cli.StringFlag{
						Name:  "query",
						Usage: "Query products using free text",
					},
					// FilterByTasteClockBody
					&RangeFlag[int]{
						Name:        "taste-clock-body",
						Usage:       "Filters products of a certain body where 0 (minimum) is a thin body and 12 (maximum) is a full body",
						HideDefault: true,
					},
					// FilterByTasteClockBitterness
					&RangeFlag[int]{
						Name:        "taste-clock-bitterness",
						Usage:       "Filter products of a certain bitterness where 0 (minimum) is not bitter at all and 12 (maximum) is very bitter",
						HideDefault: true,
					},
					// FilterByTasteClockSweetness
					&RangeFlag[int]{
						Name:        "taste-clock-sweetness",
						Usage:       "Filter products of a certain sweetness where 0 (minimum) is not sweet at all and 12 (maximum) is very sweet",
						HideDefault: true,
					},
					// FilterByTasteClockSmokiness
					&RangeFlag[int]{
						Name:        "taste-clock-smokiness",
						Usage:       "Filter products of a certain smokiness where 0 (minimum) is not smoky at all and 12 (maximum) is very smoky",
						HideDefault: true,
					},
					// FilterByVintage
					&cli.IntSliceFlag{
						Name: "vintage",
					},
					// FilterByProductLaunch
					&RangeFlag[string]{
						Name:        "product-launch",
						HideDefault: true,
					},
					// FilterByAlcoholPercentage
					&RangeFlag[int]{
						Name:        "alcohol-percentage",
						HideDefault: true,
					},
					// FilterBySugarContent
					&RangeFlag[int]{
						Name:        "sugar-content",
						HideDefault: true,
					},
					// FilterByGrapes
					&cli.StringSliceFlag{
						Name: "grapes",
					},
					// FilterByMatch
					&cli.StringSliceFlag{
						Name:  "match",
						Usage: "Filter products that fit with a taste, such as 'Aperitif', 'Asiatiskt' or 'Kött'",
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
						Name:        "volume",
						HideDefault: true,
					},
					// FilterByPackaging
					&cli.StringFlag{
						Name:  "packaging-category",
						Usage: "Filter products that use a specific packaging, such as 'Flaska'",
					},
					&cli.StringFlag{
						Name:  "packaging-subcategory",
						Usage: "Filter products that use a specific packaging, such as 'Glasflaska'",
					},
					// FilterByPrice
					&RangeFlag[int]{
						Name:        "price",
						HideDefault: true,
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
					&cli.StringFlag{
						Name:    "search",
						Aliases: []string{"q"},
						Usage:   "Optional search query",
					},
				},
			},
			{
				Name:   "stock",
				Usage:  "Get current stock",
				Action: ActionStock,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "store-id",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "product-id",
						Required: true,
					},
				},
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := root.Run(ctx, os.Args)
	cancel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
