package main

import (
	"db"
	"encoding/gob"
	"log"
	"math/rand"
	"network"
	"param"
	"phys/vect"
	"time"
)

var (
	addr = ":9696"
	s    *network.Connection
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	gob.Register(param.Player{})
	gob.Register(vect.Vect{})

	s = network.NewHandlers(map[string]network.Handler{"auth": auth})
	if err := s.Server(addr); err != nil {
		log.Fatalln(err)
	}

	for {
		log.Printf("listen '%s', clients: '%v' \n", addr, s.Clients)
		sendEnemy()

		time.Sleep(10 * time.Second)
	}
}

func auth(req *network.Request) interface{} {
	log.Println("auth", req.Data.(string))
	return db.GetPlayer(req.Data.(string))
}

func sendEnemy() {
	x := float32(rand.Intn(60) - 30)
	y := float32(rand.Intn(60) - 30)
	err := s.Broadcast("loadEnemy", "", vect.Vect{x, y})
	if err != nil {
		log.Println(err)
	}
}
