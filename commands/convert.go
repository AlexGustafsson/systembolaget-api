package commands

import (
	"fmt"
	"github.com/alexgustafsson/systembolaget-api/systembolaget"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
)

func convert(sourceName string, context *cli.Context) error {
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")
	inputExtension := strings.ToLower(filepath.Ext(input))
	outputExtension := strings.ToLower(filepath.Ext(output))

	log.Debugf("Attempting to convert source %s, %s (%s) to %s (%s)", sourceName, input, inputExtension, output, outputExtension)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var source systembolaget.Source
	if sourceName == "assortment" {
		source = &systembolaget.Assortment{}
	} else if sourceName == "inventory" {
		source = &systembolaget.Inventory{}
	} else if sourceName == "stores" {
		source = &systembolaget.Stores{}
	}

	if inputExtension == ".xml" {
		err = source.ParseFromXML(bytes)
	} else if inputExtension == ".json" {
		err = source.ParseFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputExtension)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputExtension == ".xml" {
		outputBytes, err = source.ConvertToXML(pretty)
	} else if outputExtension == ".json" {
		outputBytes, err = source.ConvertToJSON(pretty)
	} else {
		return fmt.Errorf("Unsupported output format: %s", outputExtension)
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

func convertAssortmentCommand(context *cli.Context) error {
	return convert("inventory", context)
}

func convertInventoryCommand(context *cli.Context) error {
	return convert("inventory", context)
}

func convertStoresCommand(context *cli.Context) error {
	return convert("stores", context)
}
