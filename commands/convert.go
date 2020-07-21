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

func convertAssortmentCommand(context *cli.Context) error {
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")
	inputExtension := strings.ToLower(filepath.Ext(input))
	outputExtension := strings.ToLower(filepath.Ext(output))

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputExtension, output, outputExtension)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var assortment *systembolaget.Assortment
	if inputExtension == ".xml" {
		assortment, err = systembolaget.ParseAssortmentFromXML(bytes)
	} else if inputExtension == ".json" {
		assortment, err = systembolaget.ParseAssortmentFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputExtension)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputExtension == ".xml" {
		outputBytes, err = assortment.ConvertToXML(pretty)
	} else if outputExtension == ".json" {
		outputBytes, err = assortment.ConvertToJSON(pretty)
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

func convertInventoryCommand(context *cli.Context) error {
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")
	inputExtension := strings.ToLower(filepath.Ext(input))
	outputExtension := strings.ToLower(filepath.Ext(output))

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputExtension, output, outputExtension)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var inventory *systembolaget.Inventory
	if inputExtension == ".xml" {
		inventory, err = systembolaget.ParseInventoryFromXML(bytes)
	} else if inputExtension == ".json" {
		inventory, err = systembolaget.ParseInventoryFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputExtension)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputExtension == ".xml" {
		outputBytes, err = inventory.ConvertToXML(pretty)
	} else if outputExtension == ".json" {
		outputBytes, err = inventory.ConvertToJSON(pretty)
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

func convertStoresCommand(context *cli.Context) error {
	input := context.String("input")
	output := context.String("output")
	pretty := context.Bool("pretty")
	inputExtension := strings.ToLower(filepath.Ext(input))
	outputExtension := strings.ToLower(filepath.Ext(output))

	log.Debugf("Attempting to convert %s (%s) to %s (%s)", input, inputExtension, output, outputExtension)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	var stores *systembolaget.Stores
	if inputExtension == ".xml" {
		stores, err = systembolaget.ParseStoreFromXML(bytes)
	} else if inputExtension == ".json" {
		stores, err = systembolaget.ParseStoreFromJSON(bytes)
	} else {
		return fmt.Errorf("Unsupported input format: %s", inputExtension)
	}
	if err != nil {
		return err
	}

	var outputBytes []byte
	if outputExtension == ".xml" {
		outputBytes, err = stores.ConvertToXML(pretty)
	} else if outputExtension == ".json" {
		outputBytes, err = stores.ConvertToJSON(pretty)
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
