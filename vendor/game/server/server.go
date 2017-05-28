package main

import (
	"engine"
	controllers "game/network-controllers"
	"log"
	"scenes"
)

var (
	addr = ":9696"
)

func main() {
	log.SetFlags(log.Lshortfile)
	engine.Server(server)
}

func server() {
	scenes.Load("scene00")

	controllers.NewServer(addr)
	log.Println("server listen on addr", addr)

	engine.AddCallback(controllers.SendServersState)
}
