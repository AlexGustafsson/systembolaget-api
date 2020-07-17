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

func convertAssortmentCommand(context *cli.Context) error {
	inputFormat := strings.ToLower(context.String("input-format"))
	outputFormat := strings.ToLower(context.String("output-format"))
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputFormat, output, outputFormat)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var assortment *systembolaget.Assortment
	if inputFormat == "xml" {
		assortment, err = systembolaget.ParseAssortmentFromXML(bytes)
	} else if inputFormat == "json" {
		assortment, err = systembolaget.ParseAssortmentFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputFormat)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputFormat == "xml" {
		outputBytes, err = assortment.ConvertToXML(pretty)
	} else if outputFormat == "json" {
		outputBytes, err = assortment.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", outputFormat)
	}
	if err != nil {
		return err
	}

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

func convertInventoryCommand(context *cli.Context) error {
	inputFormat := strings.ToLower(context.String("input-format"))
	outputFormat := strings.ToLower(context.String("output-format"))
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputFormat, output, outputFormat)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var inventory *systembolaget.Inventory
	if inputFormat == "xml" {
		inventory, err = systembolaget.ParseInventoryFromXML(bytes)
	} else if inputFormat == "json" {
		inventory, err = systembolaget.ParseInventoryFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputFormat)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputFormat == "xml" {
		outputBytes, err = inventory.ConvertToXML(pretty)
	} else if outputFormat == "json" {
		outputBytes, err = inventory.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", outputFormat)
	}
	if err != nil {
		return err
	}

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

func convertStoresCommand(context *cli.Context) error {
	inputFormat := strings.ToLower(context.String("input-format"))
	outputFormat := strings.ToLower(context.String("output-format"))
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputFormat, output, outputFormat)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var stores *systembolaget.Stores
	if inputFormat == "xml" {
		stores, err = systembolaget.ParseStoreFromXML(bytes)
	} else if inputFormat == "json" {
		stores, err = systembolaget.ParseStoreFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputFormat)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputFormat == "xml" {
		outputBytes, err = stores.ConvertToXML(pretty)
	} else if outputFormat == "json" {
		outputBytes, err = stores.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", outputFormat)
	}
	if err != nil {
		return err
	}

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
