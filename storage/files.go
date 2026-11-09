package storage

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

var basePath string
var compress bool

func InitFS(bp string, c bool) {
	basePath = bp
	compress = c
	log.Println("File paths: " + path("fileKey"))

	if d, err := os.Stat(bp); os.IsNotExist(err) || !d.IsDir() {
		os.MkdirAll(bp, 0750)
	}
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

func DeleteFile(key string) error {
	err := os.Remove(path(key))
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func path(key string) string {
	return basePath + string(os.PathSeparator) + key
}
