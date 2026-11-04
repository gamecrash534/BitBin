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
	conf := getConfig(os.Args[1:])

	storage.Init(conf.StorageDir, true)
	log.Println("Running on port " + strconv.Itoa(int(conf.Port)))
	web.RunServer(conf.Port, conf.Link, conf.LinkTitle)
}

func getConfig(args []string) Config {
	var conf Config
	flags := parseFlags(args)

	for k, v := range flags {
		switch k {
		case "p":
			fallthrough
		case "port":
			{
				l, err := strconv.Atoi(v)
				if err != nil || l < 0 {
					log.Fatalln("The value for the port flag must be a positive integer!")
				}

				conf.Port = uint(l)
			}

		case "d":
			conf.StorageDir = v
		case "l":
			conf.Link = v
		case "t":
			conf.LinkTitle = v
		}
	}

	return conf
}

func parseFlags(v []string) map[string]string {
	flags := map[string]string{}

	for i := range v {
		if v[i][0] != '-' {
			continue
		}
		flag := v[i][1:]

		parts := strings.SplitN(flag, "=", 2)
		if len(parts) < 2 {
			log.Fatalln("Need to provide a value to the '" + flag + "' flag!")
		} else {
			flags[parts[0]] = parts[1]
		}
	}

	return flags
}
