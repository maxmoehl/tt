package storage

import (
	"github.com/maxmoehl/tt/storage/file"
	"github.com/maxmoehl/tt/types"
)

var s types.Interface

func init() {
	var err error
	s, err = file.New()
	if err != nil {
		panic(err.Error())
	}
}
