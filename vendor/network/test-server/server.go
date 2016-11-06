package main

import (
	"log"
	"network"
	"os"
	"time"
)

var (
	addr = ":9692"
	s    *network.Connection
)

func main() {
	s = network.NewHandlers(map[string]network.Handler{"ping": ping})
	if err := s.Server(addr); err != nil {
		log.Fatalln(err)
	}

	log.Printf("server will listen on port `%s` 5 seconds", addr)
	time.Sleep(5 * time.Second)
}

func ping(req *network.Request) interface{} {
	return req.Data
}

func exit(req *network.Request) interface{} {
	s.Close()
	os.Exit(0)
	return nil
}
