package storage

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
)

var basePath string
var compress bool

func Init(bp string, c bool) {
	basePath = bp
	compress = c
	log.Println("File paths: " + path("fileKey"))
}

func StoreFile(key string, data []byte, alreadyCompressed bool) error {
	f, err := os.Create(path(key))
	if err != nil {
		return err
	}

	defer f.Close()

	if alreadyCompressed {
		_, err = f.Write(data)
		log.Println("Stored file with key '" + key + "'")
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

	log.Println("Stored file with key '" + key + "'")
	return gzw.Close()
}

func GetContents(key string) ([]byte, error) {
	if ex, err := Exists(key); err != nil {
		return nil, err
	} else if ex == false {
		log.Println("Could not find requested file with key '" + key + "'")
		return nil, nil
	}

	f, err := os.Open(path(key))
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

	log.Println("Read file with key '" + key + "'")
	return buf.Bytes(), nil
}

func Exists(key string) (bool, error) {
	_, err := os.Stat(path(key))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func path(key string) string {
	return basePath + string(os.PathSeparator) + key
}
