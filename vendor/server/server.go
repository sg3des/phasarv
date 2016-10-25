package main

import (
	"db"
	"encoding/gob"
	"log"
	"network"
	"param"
	"time"
)

var (
	addr = ":9696"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	server()
	network.AddRoute("Auth", Auth)

	for {
		log.Printf("listen on '%s'\n", addr)
		time.Sleep(time.Minute)
	}
}

func server() {
	gob.Register(param.Player{})

	if err := network.Server(addr); err != nil {
		log.Fatal(err)
	}
}

func Auth(data interface{}) interface{} {
	log.Println("Auth", data.(string))
	return db.GetPlayer(data.(string))
}
