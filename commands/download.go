package commands

import (
	"github.com/alexgustafsson/systembolaget-api/systembolaget"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strings"
)

func downloadAssortmentCommand(context *cli.Context) error {
	response, err := systembolaget.DownloadAssortment()
	if err != nil {
		return err
	}

	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	if format == "xml" {
		xml, err := response.ConvertToXML()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")

		if output == "" {
			os.Stdout.Write(xml)
		} else {
			ioutil.WriteFile(output, xml, 0644)
		}
	} else if format == "json" {
		json, err := response.ConvertToJSON()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")
		if output == "" {
			os.Stdout.Write(json)
		} else {
			ioutil.WriteFile(output, json, 0644)
		}
	}

	return nil
}

func downloadInventoryCommand(context *cli.Context) error {
	response, err := systembolaget.DownloadInventory()
	if err != nil {
		return err
	}

	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	if format == "xml" {
		xml, err := response.ConvertToXML()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")
		if output == "" {
			os.Stdout.Write(xml)
		} else {
			ioutil.WriteFile(output, xml, 0644)
		}
	} else if format == "json" {
		json, err := response.ConvertToJSON()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")
		if output == "" {
			os.Stdout.Write(json)
		} else {
			ioutil.WriteFile(output, json, 0644)
		}
	}

	return nil
}

func downloadStoresCommand(context *cli.Context) error {
	response, err := systembolaget.DownloadStores()
	if err != nil {
		return err
	}

	format := strings.ToLower(context.String("format"))
	output := context.String("output")
	if format == "xml" {
		xml, err := response.ConvertToXML()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")
		if output == "" {
			os.Stdout.Write(xml)
		} else {
			ioutil.WriteFile(output, xml, 0644)
		}
	} else if format == "json" {
		json, err := response.ConvertToJSON()
		if err != nil {
			return err
		}

		log.Debug("Successfully processed assortment API")
		if output == "" {
			os.Stdout.Write(json)
		} else {
			ioutil.WriteFile(output, json, 0644)
		}
	}

	return nil
}
