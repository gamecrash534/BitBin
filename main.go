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
	Port       int
	Link       string
	LinkTitle  string
	FileTtlMin int64
}

var cleanupTask *time.Ticker

func main() {
	conf := getConfig()

	storage.Init(conf.StorageDir, storage.Config{
		CompressFiles:     true,
		DbCompressOnStart: true,
		DbCacheMiB:        2,
		DataTtlMin:        conf.FileTtlMin,
	})

	defer func() {
		cleanupTask.Stop()
		storage.Close()
	}()

	log.Println("Running on port " + strconv.Itoa(conf.Port))

	go createCleanupTask()
	web.RunServer(conf.Port, conf.Link, conf.LinkTitle)
}

func getConfig() Config {
	conf := Config{
		StorageDir: "./data",
		Port:       8080,
		Link:       "https://git.gamecrash.dev/game.crash/BitBin",
		LinkTitle:  "source",
		FileTtlMin: 43200,
	}

	if d := getenv("directory"); d != nil {
		conf.StorageDir = *d
	}

	if port := getenvInt("port"); port != nil {
		conf.Port = *port
	}

	if l := getenv("link"); l != nil {
		conf.Link = *l
	}

	if lt := getenv("link_title"); lt != nil {
		conf.LinkTitle = *lt
	}

	if ttl := getenvInt("paste_ttl"); ttl != nil {
		conf.FileTtlMin = int64(*ttl)
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

func getenvInt(key string) *int {
	v := getenv(key)
	if v == nil {
		return nil
	}

	c, err := strconv.Atoi(*v)
	if err != nil {
		return nil
	}

	return &c
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
