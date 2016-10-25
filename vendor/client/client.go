package client

import (
	"encoding/gob"
	"log"
	"network"
	"param"
)

func Connect(addr string) {
	gob.Register(param.Player{})

	if err := network.Client(addr); err != nil {
		log.Fatalln("failed connect to server, error:", err)
	}

}

func AddRoutes(funcs map[string]func(interface{}) interface{}) {
	for funcname, f := range funcs {
		network.AddRoute(funcname, f)
	}
}

func Authorize(name string) {
	log.Println(network.Connection.Routes)
	if err := network.SendMessage("Auth", "CreateLocalPlayer", name); err != nil {
		log.Fatalln("failed authorize", err)
	}
}
