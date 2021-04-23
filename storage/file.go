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

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"

	"github.com/google/uuid"
)

type file struct {
	timers types.Timers
}

func (f *file) GetTimer(uuid uuid.UUID) (types.Timer, error) {
	for _, t := range f.timers {
		if t.Uuid == uuid {
			return t, nil
		}
	}
	return types.Timer{}, fmt.Errorf("no timer found for uuid %s", uuid.String())
}

func (f *file) GetRunningTimer() (types.Timer, error) {
	for _, t := range f.timers {
		if t.End.IsZero() {
			return t, nil
		}
	}
	return types.Timer{}, fmt.Errorf("no running timer found")
}

func (f *file) GetLastTimer(running bool) (types.Timer, error) {
	return f.timers.Last(running), nil
}

func (f *file) GetTimers(filter types.Filter) (types.Timers, error) {
	return f.timers.Filter(filter), nil
}

func (f *file) RunningTimerExists() (bool, error) {
	for _, t := range f.timers {
		if t.Running() {
			return true, nil
		}
	}
	return false, nil
}

func (f *file) StoreTimer(newTimer types.Timer) error {
	exists := false
	for _, t := range f.timers {
		if t.Uuid == newTimer.Uuid {
			exists = true
			break
		}
	}
	if exists {
		return fmt.Errorf("timer with uuid %s already exists", newTimer.Uuid.String())
	}
	f.timers = append(f.timers, newTimer)
	return f.write()
}

func (f *file) UpdateTimer(updatedTimer types.Timer) error {
	for i, t := range f.timers {
		if t.Uuid == updatedTimer.Uuid {
			f.timers[i].End = updatedTimer.End
			break
		}
	}
	return f.write()
}

func (f *file) write() error {
	fileWriter, err := os.Create(getStorageFile())
	if err != nil {
		return err
	}
	return json.NewEncoder(fileWriter).Encode(f.timers)
}

// NewFile initializes and returns a new storage interface that can be used
// to access data.
func NewFile() (types.Storage, error) {
	var f file
	if !storageFileExists() {
		return &f, nil
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
