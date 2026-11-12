package main

import (
	"bitbin/storage"
	"bitbin/web"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	StorageDir string
	Port       uint
	Link       string
	LinkTitle  string
}

var cleanupTask *time.Ticker

func main() {
	conf := getConfig()

	storage.Init(conf.StorageDir, storage.Config{
		CompressFiles:     true,
		DbCompressOnStart: true,
		DbCacheMiB:        2,
		DataTtlMin:        43200,
	})

	defer func() {
		cleanupTask.Stop()
		storage.Close()
	}()

	log.Println("Running on port " + strconv.Itoa(int(conf.Port)))

	go createCleanupTask()
	web.RunServer(conf.Port, conf.Link, conf.LinkTitle)
}

func getConfig() Config {
	conf := Config{
		StorageDir: "./data",
		Port:       8080,
		Link:       "https://git.gamecrash.dev/game.crash/BitBin",
		LinkTitle:  "source",
	}

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

func createCleanupTask() {
	cleanupTask = time.NewTicker(1 * time.Minute)

	done := make(chan bool)
	for {
		select {
		case <-done:
			break
		case _ = <-cleanupTask.C:
			err := storage.Cleanup()
			if err != nil {
				log.Println("Error while performing cleanup:", err)
				log.Println("Will stop performing cleanups for now")

				done <- true
				cleanupTask.Stop()
			}
		}
	}
}
