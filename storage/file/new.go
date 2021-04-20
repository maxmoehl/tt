package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"
)

func New() (types.Interface, error) {
	var f file
	if !storageFileExists() {
		return &file{}, nil
	}
	fileReader, err := os.Open(getStorageFile())
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(fileReader).Decode(&f.timers)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func getStorageFile() string {
	return filepath.Join(config.HomeDir(), "storage.json")
}

func storageFileExists() bool {
	_, err := os.Stat(getStorageFile())
	if os.IsNotExist(err) {
		return false
	}
	return true
}
