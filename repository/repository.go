package repository

import (
	"github.com/sdomino/scribble"
)

var (
	db *scribble.Driver
)

func init() {
	var err error
	db, err = scribble.New(dataId, nil)
	if err != nil {
		panic(err)
	}
}