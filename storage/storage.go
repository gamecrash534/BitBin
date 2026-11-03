package storage

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
)

var basePath string
var compress bool

func Init(bp string, c bool) {
	basePath = bp
	compress = c
}

func StoreFile(key string, data []byte) error {
	f, err := os.Create(basePath + "/" + key)
	if err != nil {
		return err
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	var compression int
	if compress {
		compression = gzip.DefaultCompression
	} else {
		compression = gzip.NoCompression
	}

	gzw, err := gzip.NewWriterLevel(f, compression)
	if err != nil {
		return err
	}
	defer gzw.Close()

	gzw.Write(data)

	return err
}

func GetContents(key string) ([]byte, error) {
	if ex, err := Exists(key); err != nil {
		return nil, err
	} else if ex == false {
		return nil, nil
	}

	f, err := os.Open(basePath + "\\" + key)
	if err != nil {
		return nil, err
	}

	gzw, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzw.Close()

	buf := bytes.NewBuffer(make([]byte, 0))

	io.Copy(buf, gzw)

	return buf.Bytes(), nil
}

func Exists(key string) (bool, error) {
	_, err := os.Stat(basePath + "\\" + key)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
