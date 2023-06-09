package repository

import (
	"fmt"

	"github.com/sdomino/scribble"
)

var (
	db *scribble.Driver
)

func init() {
	var err error
	db, err = scribble.New("data", nil)
	if err != nil {
		panic(err)
	}
}

func formatKey(collectionPrefix, entity string) string {
	return fmt.Sprintf("%s-%s", collectionPrefix, entity)
}
