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

package tt

import (
	"github.com/maxmoehl/tt/config"
)

var s Storage

func init() {
	err := initStorage()
	if err != nil {
		panic(err.Error())
	}
}

func initStorage() error {
	var err error
	switch config.Get().StorageType {
	case config.StorageTypeFile:
		s, err = NewFile()
	case config.StorageTypeSQLite:
		s, err = NewSQLite()
	}
	if err != nil {
		return err
	}
	return nil
}
