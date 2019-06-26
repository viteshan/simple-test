package db

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestVersion(t *testing.T) {
	d, err := leveldb.OpenFile("testdata-consensus", nil)
	if err != nil {
		panic(err)
	}

}
