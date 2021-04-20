/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"
)

// New initializes and returns a new storage interface that can be used
// to access data.
func New() (types.Storage, error) {
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
