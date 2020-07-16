package commands

import (
	"fmt"
	"github.com/alexgustafsson/systembolaget-api/systembolaget"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strings"
)

func downloadAssortmentCommand(context *cli.Context) error {
	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debug("Attempting to download assortment")
	response, err := systembolaget.DownloadAssortment()
	if err != nil {
		return err
	}

	log.Debug("Downloaded data, converting to output format")
	var outputBytes []byte
	if format == "xml" {
		outputBytes, err = response.ConvertToXML(pretty)
	} else if format == "json" {
		outputBytes, err = response.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", format)
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

func downloadInventoryCommand(context *cli.Context) error {
	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debug("Attempting to download inventory")
	response, err := systembolaget.DownloadInventory()
	if err != nil {
		return err
	}

	log.Debug("Downloaded data, converting to output format")
	var outputBytes []byte
	if format == "xml" {
		outputBytes, err = response.ConvertToXML(pretty)
	} else if format == "json" {
		outputBytes, err = response.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", format)
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

func downloadStoresCommand(context *cli.Context) error {
	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debug("Attempting to download stores")
	response, err := systembolaget.DownloadStores()
	if err != nil {
		return err
	}

	log.Debug("Downloaded data, converting to output format")
	var outputBytes []byte
	if format == "xml" {
		outputBytes, err = response.ConvertToXML(pretty)
	} else if format == "json" {
		outputBytes, err = response.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", format)
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
