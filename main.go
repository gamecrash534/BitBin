package main

import (
	"bitbin/storage"
	"bitbin/web"
	"log"
)

func main() {
	storage.Init("/dev/null", true)
	log.Println("running")
	web.RunServer(8080, "https://github.com/gamecrash534", "github")
}
