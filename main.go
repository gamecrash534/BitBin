package main

import (
	"bitbin/storage"
	"bitbin/web"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	StorageDir string
	Port       uint
	Link       string
	LinkTitle  string
}

func main() {
	conf := getConfig()

	storage.Init(conf.StorageDir, true)
	log.Println("Running on port " + strconv.Itoa(int(conf.Port)))
	web.RunServer(conf.Port, conf.Link, conf.LinkTitle)
}

func getConfig() Config {
	var conf Config

	if d := getenv("directory"); d != nil {
		conf.StorageDir = *d
	}

	if port := getenvUint("port"); port != nil {
		conf.Port = *port
	}

	if l := getenv("link"); l != nil {
		conf.Link = *l
	}

	if lt := getenv("link_title"); lt != nil {
		conf.LinkTitle = *lt
	}

	return conf
}

func getenv(key string) *string {
	v := os.Getenv("BITBIN_" + strings.ToUpper(key))

	if v == "" {
		return nil
	}

	return &v
}

func getenvUint(key string) *uint {
	v := getenv(key)
	if v == nil {
		return nil
	}

	c, err := strconv.Atoi(*v)
	if err != nil {
		return nil
	}

	return new(uint(c))
}
