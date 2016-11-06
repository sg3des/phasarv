package main

import (
	"encoding/gob"
	"log"
	"network"
	"param"
	"phys/vect"
)

var (
	c *network.Connection
)

func Connect(addr string) {
	gob.Register(param.Player{})
	gob.Register(vect.Vect{})

	c = network.NewHandlers(map[string]network.Handler{
		"loadLocalPlayer": loadLocalPlayer,
		"loadEnemy":       loadEnemy,
	})
	if err := c.Client(addr); err != nil {
		log.Fatalln(err)
	}

}

func Authorize(name string) {
	if err := c.SendMessage("auth", "loadLocalPlayer", name); err != nil {
		log.Fatalln("failed authorize", err)
	}
}

func loadLocalPlayer(req *network.Request) interface{} {
	CreateLocalPlayer(req.Data.(param.Player))
	return nil
}

func loadEnemy(req *network.Request) interface{} {
	pos := req.Data.(vect.Vect)
	log.Println("load enemy", pos)
	createEnemy(pos.X, pos.Y)
	return nil
}
