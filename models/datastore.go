package models

import (
	"encoding/gob"
	"os"

	"github.com/antsy/tournament/utils"
)

/**
 * Save
 */
func Persist() error {
	file, err := os.Create(utils.StoragePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(Games)
	}
	file.Close()
	return err
}

/**
 * Load
 */
func Retrive() error {
	var data = new([]Tournament)
	file, err := os.Open(utils.StoragePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(data)
	}
	file.Close()
	Games = *data
	return err
}
