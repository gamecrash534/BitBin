package main

import (
	"bitbin/web"
	"log"
)

func main() {
	log.Println("running")
	web.RunServer(8080, "https://github.com/gamecrash534", "github")
}
