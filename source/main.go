package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createOutputFolders() {
	outputPath := filepath.Join(".", "output/xml")
	os.MkdirAll(outputPath, os.ModePerm)

	outputPath = filepath.Join(".", "output/json")
	os.MkdirAll(outputPath, os.ModePerm)
}

func main() {
	var storesAPI, assortmentAPI, inventoryAPI API

	storesAPI = &StoresAPI{}
	assortmentAPI = &AssortmentAPI{}
	inventoryAPI = &InventoryAPI{}

	createOutputFolders()

	err := storesAPI.Process()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while processing stores API: %v", err))
	} else {
		fmt.Println("Successfully processed stores API")
	}

	err = assortmentAPI.Process()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while processing assortment API: %v", err))
	} else {
		fmt.Println("Successfully processed assortment API")
	}

	err = inventoryAPI.Process()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while processing inventory API: %v", err))
	} else {
		fmt.Println("Successfully processed inventory API")
	}
}
