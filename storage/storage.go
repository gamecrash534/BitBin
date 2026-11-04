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

func StoreFile(key string, data []byte, alreadyCompressed bool) error {
	f, err := os.Create(basePath + "/" + key)
	if err != nil {
		return err
	}

	defer f.Close()

	if alreadyCompressed {
		_, err = f.Write(data)
		return err
	}

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

	if _, err = gzw.Write(data); err != nil {
		gzw.Close()
		return err
	}

	return gzw.Close()
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
	defer f.Close()

	gzw, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzw.Close()

	buf := bytes.NewBuffer(make([]byte, 0))

	if _, err = io.Copy(buf, gzw); err != nil {
		return nil, err
	}

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
