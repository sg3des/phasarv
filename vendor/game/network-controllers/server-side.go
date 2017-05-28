package controllers

import (
	"game"
	"game/db"
	"log"
	"math/rand"
	"net"
	"network"
	"phys/vect"
	"time"
)

var (
	s       Server
	clients = make(map[string]*cliPlayer)
)

func NewServer(addr string) {
	conn = network.NewConnection(s)
	err := conn.Server(addr)
	if err != nil {
		log.Fatalln(err)
	}

}

func SendServersState(float32) bool {
	if len(clients) == 0 {
		return true
	}

	timeNow := time.Now()

	var states ServersState
	for _, c := range clients {
		// log.Println(c.p.Name, c.p.CurrParam.Health, timeNow.After(c.deadline))
		if timeNow.After(c.deadline) {
			conn.DeleteClient(c.addr)
			conn.Broadcast(Client.RemovePlayer, nil, c.p.Name)
			delCliPlayer(c.p.Name)
			c = nil
			continue
		}
		states = append(states, c.GetServerState())
	}

	err := conn.Broadcast(Client.PlayersServerState, nil, states)
	if err != nil {
		log.Println(err)
	}

	return true
}

func resetDeadline() time.Time {
	return time.Now().Add(durationDeadline)
}

func sendEnemy() {
	x := float32(rand.Intn(60) - 30)
	y := float32(rand.Intn(60) - 30)
	err := conn.Broadcast("loadEnemy", "", vect.Vect{x, y})
	if err != nil {
		log.Println(err)
	}
}

func newCliPlayer(name string, req *network.Request) (c *cliPlayer) {
	var addr *net.UDPAddr
	if req != nil {
		addr = req.RemoteAddr
	}

	c = &cliPlayer{db.GetPlayer(name), addr, resetDeadline()}
	clients[name] = c

	// p = db.GetPlayer(name)

	// p := db.GetPlayer(name)

	return c
}

func addCliPlayer(p *game.Player) *cliPlayer {
	c := &cliPlayer{p, nil, time.Time{}}
	clients[p.Name] = c
	return c
}

func delCliPlayer(name string) {
	delete(clients, name)
}

//delPlayerByString remove player by string key
func delCliPlayerByReq(req *network.Request) {
	addr := req.RemoteAddr.String()
	for _, c := range clients {
		if c.addr.String() == addr {
			delete(clients, c.p.Name)
			return
		}
	}
}

func lookupCliPlayer(name string, req *network.Request) (c *cliPlayer, ok bool) {

	if req != nil {
		addr := req.RemoteAddr.String()
		for _, c := range clients {
			if c.addr.String() == addr {
				return c, true
			}
		}
	}

	c, ok = clients[name]

	return
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

	c := newCliPlayer(name, req)
	game.CreatePlayer(c.p)
	// p.CreatePlayer()
	c.p.Object.AddCallback(c.p.ClientCursor, c.p.Movement, c.p.PlayerRotation)

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

	c, ok := lookupCliPlayer("", req)
	// p, ok := players[req.RemoteAddr.String()]
	if !ok {
		log.Println("WARNING player is not connected", req.RemoteAddr.String(), clients)
		return nil
	}

	c.deadline = resetDeadline()

	s.UpdatePlayer(c.p)

	return ServerState{
		Vel:  c.p.Object.Velocity(),
		AVel: c.p.Object.AngularVelocity(),
		Pos:  c.p.Object.PositionVect(),
		Rot:  c.p.Object.Rotation(),
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
