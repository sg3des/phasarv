package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"network"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	gob.Register(Data{})

	if len(os.Args) < 2 {
		log.Fatalln("need set mode - server or client how argument")
	}

	addr := "127.0.0.1:9696"
	if os.Args[1] == "server" {
		server(addr)
	} else {
		client(addr)
	}

}

func server(addr string) {
	var h Server
	conn := network.NewConnection(h)
	err := conn.Server(addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	for {
		time.Sleep(2 * time.Second)
		log.Println("clients:", len(conn.Clients))
	}
}

func client(addr string) {
	var h Client
	conn := network.NewConnection(h)
	err := conn.Client(addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.SendMessage(Server.Hello, Client.Hello, Data{"hello i`m client"})
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(1 * time.Second)
	conn.Close()
}

type Data struct {
	Msg string
}

type Server struct{}
type Client struct{}

func (Server) Hello(req *network.Request) interface{} {
	s, ok := req.Data.(Data)
	if !ok {
		return nil
	}
	fmt.Printf("server get message: %s, from %s\n", s.Msg, req.RemoteAddr)
	return Data{"hello i`m server"}
}

func (Client) Hello(req *network.Request) interface{} {
	s, ok := req.Data.(Data)
	if !ok {
		return nil
	}
	fmt.Println("client get message:", s.Msg)
	return nil
}
