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
