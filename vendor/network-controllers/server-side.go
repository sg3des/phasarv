package controllers

import (
	"db"
	"game"
	"log"
	"math/rand"
	"network"
	"phys/vect"
)

var (
	s       Server
	players = make(map[string]*game.Player)
)

func NewServer(addr string) {
	conn = network.NewConnection(s)
	err := conn.Server(addr)
	if err != nil {
		log.Fatalln(err)
	}
}

func SendServersState(float32) bool {
	if len(players) == 0 {
		return true
	}

	var states ServersState
	for _, p := range players {
		states = append(states, GetServerState(p))
	}
	err := conn.Broadcast(Client.PlayersServerState, nil, states)
	if err != nil {
		log.Println(err)
	}

	return true
}

func sendEnemy() {
	x := float32(rand.Intn(60) - 30)
	y := float32(rand.Intn(60) - 30)
	err := conn.Broadcast("loadEnemy", "", vect.Vect{x, y})
	if err != nil {
		log.Println(err)
	}
}

func addPlayer(req *network.Request, name string) *game.Player {
	p := db.GetPlayer(name)
	players[req.RemoteAddr.String()] = p
	return p
}

func delPlayer(req *network.Request) {
	delete(players, req.RemoteAddr.String())
}

func lookupPlayer(req *network.Request) (*game.Player, bool) {
	p, ok := players[req.RemoteAddr.String()]
	return p, ok
}

//
//
// Server side controllers handlers
type Server struct{}

func (Server) Authorize(req *network.Request) interface{} {
	name, ok := req.Data.(string)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
	}
	log.Println("auth", name)

	p := addPlayer(req, name)
	p.CreatePlayer()
	p.Object.AddCallback(p.ClientCursor, p.Movement, p.PlayerRotation)

	return db.GetPlayer(name)
}

func (Server) GetPlayer(req *network.Request) interface{} {
	name, ok := req.Data.(string)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
	}

	return db.GetPlayer(name)
}

func (Server) PlayerState(req *network.Request) interface{} {
	s, ok := req.Data.(ClientState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	p, ok := players[req.RemoteAddr.String()]
	if !ok {
		log.Println("WARNING player is not connected", req.RemoteAddr.String(), players)
		return nil
	}

	s.UpdatePlayer(p)

	// log.Println(p.Object.PositionVect(), p.Cursor.PositionVect())

	return ServerState{
		Vel:  p.Object.Velocity(),
		AVel: p.Object.AngularVelocity(),
		Pos:  p.Object.PositionVect(),
		Rot:  p.Object.Rotation(),
	}
}

// func handlerPlayerState(req *network.Request) interface{} {
// 	name, ok := req.Data.(string)
// 	if !ok || name == "" {
// 		log.Println("WARNING! recieve data is not correct")
// 		return nil
// 	}

// 	p := lookupPlayer(req)
// }
