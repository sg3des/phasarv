package main

import (
	"db"
	"encoding/gob"
	"engine"
	"game"
	"log"
	"math/rand"
	"network"
	"phys/vect"
	"scene"
)

var (
	addr = ":9696"
	s    *network.Connection
)

func init() {
	log.SetFlags(log.Lshortfile)
	players = make(map[string]*game.Player)
}

func main() {
	engine.Server(server)
}

func server() {
	scene.Load("scene00")

	gob.Register(game.Player{})
	gob.Register(vect.Vect{})
	gob.Register(game.NetPacket{})

	s = network.NewHandlers(map[string]network.Handler{
		"auth":          auth,
		"playersCursor": playersCursor,
	})

	if err := s.Server(addr); err != nil {
		log.Fatalln(err)
	}

	log.Println("server listen on addr", addr)

	// for {
	// 	log.Printf("listen '%s', clients: '%v' \n", addr, s.Clients)
	// 	sendEnemy()

	// 	time.Sleep(10 * time.Second)
	// }
}

func auth(req *network.Request) interface{} {
	log.Println("auth", req.Data.(string))

	name := req.Data.(string)
	p := db.GetPlayer(name)
	p.CreatePlayer()
	p.Object.AddCallback(p.Movement, p.PlayerRotation)
	players[req.RemoteAddr.String()] = p

	return db.GetPlayer(name)
}

func sendEnemy() {
	x := float32(rand.Intn(60) - 30)
	y := float32(rand.Intn(60) - 30)
	err := s.Broadcast("loadEnemy", "", vect.Vect{x, y})
	if err != nil {
		log.Println(err)
	}
}
