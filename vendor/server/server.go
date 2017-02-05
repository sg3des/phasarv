package main

import (
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

	game.RegisterNetworkTypes()

	s = network.NewHandlers(map[string]network.Handler{
		"auth":        auth,
		"clientState": clientState,
		"getPlayer":   getPlayer,
	})

	if err := s.Server(addr); err != nil {
		log.Fatalln(err)
	}

	log.Println("server listen on addr", addr)

	engine.AddCallback(sendServersState)

	// for {
	// 	log.Printf("listen '%s', clients: '%v' \n", addr, s.Clients)
	// 	sendEnemy()

	// 	time.Sleep(10 * time.Second)
	// }
}

func sendEnemy() {
	x := float32(rand.Intn(60) - 30)
	y := float32(rand.Intn(60) - 30)
	err := s.Broadcast("loadEnemy", "", vect.Vect{x, y})
	if err != nil {
		log.Println(err)
	}
}
