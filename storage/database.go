package storage

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Data struct {
	Deletion int64
	MimeType string
}

var db *leveldb.DB

func InitDB(p string, compact bool, cache int) {
	log.Println("Init DB")
	var err error
	db, err = leveldb.OpenFile(p, &opt.Options{
		BlockCacheCapacity: cache * opt.MiB,
	})

	if compact {
		log.Println("Running DB compact")
		err := db.CompactRange(util.Range{
			Start: nil,
			Limit: nil,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB init done")
}

func Close() {
	db.Close()
}

func Exists(k string) (bool, error) {
	b, err := db.Has(b(k), nil)

	if err != nil {
		return false, err
	}

	if b {
		return true, nil
	}

	log.Println("Not found: '", k, "'")
	return false, nil
}

func PutDB(k string, d Data) error {
	r, err := ser(d)
	if err != nil {
		return err
	}

	if err := db.Put(b(k), r, nil); err != nil {
		return err
	}

	return nil
}

func GetDB(k string) (*Data, error) {
	d, err := db.Get(b(k), nil)

	if err != nil {
		return nil, err
	}

	r, err := deser(d)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func RemoveDB(k string) error {
	return db.Delete(b(k), nil)
}

func b(s string) []byte {
	return []byte(s)
}

func ser(d Data) ([]byte, error) {
	return json.Marshal(d)
}

func deser(d []byte) (*Data, error) {
	var t Data

	if err := json.Unmarshal(d, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
