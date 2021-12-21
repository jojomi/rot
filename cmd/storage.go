package cmd

import (
	"encoding/json"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

var storageFilename, _ = homedir.Expand("~/rot.data")

// load loads the rot data from the current user's home directory (the file is called rot.data)
func load() ([]RotItem, error) {
	var rotItems []RotItem

	if exists, _ := afero.Exists(AppFs, storageFilename); !exists {
		return rotItems, nil
	}

	file, err := afero.ReadFile(AppFs, storageFilename)
	if err != nil {
		return rotItems, err
	}

	json.Unmarshal(file, &rotItems)
	return rotItems, nil
}

// save saves the rot data given to the current user's home directory (the file is called rot.data)
func save(items []RotItem) error {
	jsonData, _ := json.Marshal(items)
	err := afero.WriteFile(AppFs, storageFilename, jsonData, 0640)
	return err
}
