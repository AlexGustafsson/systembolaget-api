package commands

import (
	"github.com/urfave/cli/v2"
)

var downloadFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "output, o",
		Usage: "Output file",
	},
	&cli.StringFlag{
		Name:  "format, f",
		Usage: "Output format. Eiter JSON or XML",
		Value: "JSON",
	},
}

// Commands contains all commands of the application
var Commands = []*cli.Command{
	{
		Name:   "version",
		Usage:  "Show the application's version",
		Action: versionCommand,
	},
	{
		Name:  "download",
		Usage: "Download API data",
		Subcommands: []*cli.Command{
			{
				Name:   "assortment",
				Usage:  "Download assortment data",
				Action: downloadAssortmentCommand,
				Flags:  downloadFlags,
			},
			{
				Name:   "inventory",
				Usage:  "Download inventory data",
				Action: downloadInventoryCommand,
				Flags:  downloadFlags,
			},
			{
				Name:   "stores",
				Usage:  "Download stores data",
				Action: downloadStoresCommand,
				Flags:  downloadFlags,
			},
		},
	},
}
