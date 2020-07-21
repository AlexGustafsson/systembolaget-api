package commands

import (
	"fmt"
	"github.com/alexgustafsson/systembolaget-api/systembolaget"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"os"
	"strings"
)

func download(source string, context *cli.Context) error {
	output := context.String("output")
	pretty := context.Bool("pretty")
	extension := strings.ToLower(filepath.Ext(output))

	log.Debug("Attempting to download assortment")
	var response systembolaget.Source
	var err error
	if source == "assortment" {
		response, err = systembolaget.DownloadAssortment()
	} else if source == "inventory" {
		response, err = systembolaget.DownloadInventory()
	} else if source == "stores" {
		response, err = systembolaget.DownloadStores()
	}
	if err != nil {
		return err
	}

	log.Debug("Downloaded data, converting to output format")
	var outputBytes []byte
	if extension == ".xml" {
		outputBytes, err = response.ConvertToXML(pretty)
	} else if extension == ".json" {
		outputBytes, err = response.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output extension: %s", extension)
	}
	if err != nil {
		return err
	}

	log.Debug("Converted data, writing to target")
	if output == "" {
		os.Stdout.Write(outputBytes)
	} else {
		err = ioutil.WriteFile(output, outputBytes, 0644)
	}
	if err != nil {
		return err
	}

	return nil
}

func downloadAssortmentCommand(context *cli.Context) error {
	return download("assortment", context)
}

func downloadInventoryCommand(context *cli.Context) error {
	return download("inventory", context)
}

func downloadStoresCommand(context *cli.Context) error {
	return download("stores", context)
}
