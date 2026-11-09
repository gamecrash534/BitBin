package storage

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Data struct {
	deletion int64
	mimeType string
}

var buf = bytes.NewBuffer(make([]byte, 0))
var enc = gob.NewEncoder(buf)
var dec = gob.NewDecoder(buf)

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

	return b, nil
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
	buf.Reset()

	if err := enc.Encode(d); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deser(d []byte) (*Data, error) {
	buf.Reset()
	buf.Write(d)
	var r Data

	if err := dec.Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
