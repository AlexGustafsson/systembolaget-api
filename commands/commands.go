package commands

import (
	"github.com/AlexGustafsson/systembolaget-api/commands/v1"
	"github.com/urfave/cli/v2"
)

var outputFlag = &cli.StringFlag{
	Name:    "output",
	Aliases: []string{"o"},
	Usage:   "Output file",
}

var outputFormat = &cli.StringFlag{
	Name:    "format",
	Aliases: []string{"f"},
	Usage:   "Output format",
}

var keyFlag = &cli.StringFlag{
	Name: "key",
	Usage: "API key",
}

// Commands contains all commands of the application
var Commands = []*cli.Command{
	{
		Name:   "version",
		Usage:  "Show the application's version",
		Action: versionCommand,
	},
	{
		Name:  "v1",
		Usage: "API version 1",
		Subcommands: []*cli.Command{
			{
				Name:  "product",
				Usage: "Product commands",
				Subcommands: []*cli.Command{
					{
						Name:   "fetch",
						Usage:  "Fetch a single product",
						Action: v1.ProductFetchCommand,
						Flags:  []cli.Flag{
							outputFlag,
							outputFormatFlag,
							keyFlag,
							&cli.StringFlag{
								Name:    "id",
								Usage:   "Product id",
							},
						},
					},
				},
			},
		},
	},
}
