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
		Usage: "Output format. Either JSON or XML",
		Value: "JSON",
	},
	&cli.BoolFlag{
		Name:  "pretty",
		Usage: "Pretty print output",
		Value: false,
	},
}

var convertFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "input, i",
		Usage: "Input file",
	},
	&cli.StringFlag{
		Name:  "output, o",
		Usage: "Output file",
	},
	&cli.StringFlag{
		Name:  "input-format",
		Usage: "Input format. Either JSON or XML",
		Value: "JSON",
	},
	&cli.StringFlag{
		Name:  "output-format",
		Usage: "Output format. Either JSON or XML",
		Value: "JSON",
	},
	&cli.BoolFlag{
		Name:  "pretty",
		Usage: "Pretty print output",
		Value: false,
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
	{
		Name:  "convert",
		Usage: "Convert API data from one format to another",
		Subcommands: []*cli.Command{
			{
				Name:   "assortment",
				Usage:  "Convert assortment data",
				Action: convertAssortmentCommand,
				Flags:  convertFlags,
			},
			{
				Name:   "inventory",
				Usage:  "Convert inventory data",
				Action: convertInventoryCommand,
				Flags:  convertFlags,
			},
			{
				Name:   "stores",
				Usage:  "Convert stores data",
				Action: convertStoresCommand,
				Flags:  convertFlags,
			},
		},
	},
}
