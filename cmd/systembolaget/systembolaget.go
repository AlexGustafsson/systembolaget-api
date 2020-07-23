package main

import (
	"github.com/alexgustafsson/systembolaget-api/commands"
	"github.com/alexgustafsson/systembolaget-api/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"sort"
)

var appHelpTemplate = `Usage: {{.Name}} [global options] command [command options] [arguments]

{{.Usage}}

Version: {{.Version}}

Options:
  {{range .Flags}}{{.}}
  {{end}}
Commands:
  {{range .Commands}}{{.Name}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} help command' for more information on a command.
`

var commandHelpTemplate = `Usage: tunnel {{.Name}} [options] [arguments]

{{.Usage}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .Flags}}

Options:{{range .Flags}}
   {{.}}{{end}}{{end}}
`

func setDebugOutputLevel() {
	for _, flag := range os.Args {
		if flag == "-v" || flag == "--verbose" {
			log.SetLevel(log.DebugLevel)
		}
	}
}

func commandNotFound(context *cli.Context, command string) {
	log.Errorf(
		"%s: '%s' is not a %s command. See '%s help'.",
		context.App.Name,
		command,
		context.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}

func main() {
	setDebugOutputLevel()

	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "Interact with Systembolaget's APIs."
	app.Version = version.FullVersion()
	app.CommandNotFound = commandNotFound
	app.EnableBashCompletion = true
	app.Commands = commands.Commands
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Enable verbose logging",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
