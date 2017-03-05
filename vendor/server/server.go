package main

import (
	"engine"
	"log"
	"scene"

	controllers "network-controllers"
)

var (
	addr = ":9696"
)

func main() {
	log.SetFlags(log.Lshortfile)
	engine.Server(server)
}

func server() {
	scene.Load("scene00")

	controllers.NewServer(addr)
	log.Println("server listen on addr", addr)

	engine.AddCallback(controllers.SendServersState)
}
