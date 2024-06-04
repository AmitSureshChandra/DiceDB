package main

import (
	"dicedb/config"
	"dicedb/server"
	"flag"
	"log"
)

func main() {
	setUpConfig()
	err := server.RunAsyncServer()
	if err != nil {
		log.Print(err.Error())
	}
}

func setUpConfig() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "hostname")
	flag.IntVar(&config.Port, "port", 6379, "port")
}
