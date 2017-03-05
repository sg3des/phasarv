package controllers

import (
	"db"
	"game"
	"log"
	"math/rand"
	"net"
	"network"
	"phys/vect"
	"time"
)

var (
	s       Server
	players = make(map[string]*player)
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

	timeNow := time.Now()

	var states ServersState
	for _, p := range players {
		log.Println(p.Name, timeNow.After(p.deadline))
		if timeNow.After(p.deadline) {
			conn.DeleteClient(p.addr)
			conn.Broadcast(Client.RemovePlayer, nil, p.Name)
			delPlayer(p.Name)
			p = nil
			continue
		}
		states = append(states, p.GetServerState())
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

func newPlayer(name string, req *network.Request) (p *player) {
	var addr *net.UDPAddr
	if req != nil {
		addr = req.RemoteAddr
	}

	p = &player{db.GetPlayer(name), addr, resetDeadline()}
	players[name] = p

	// p = db.GetPlayer(name)

	// p := db.GetPlayer(name)

	return p
}

func addPlayer(p *game.Player) *player {
	player := &player{p, nil, time.Time{}}
	players[p.Name] = player
	return player
}

func delPlayer(name string) {
	delete(players, name)
}

//delPlayerByString remove player by string key
func delPlayerByReq(req *network.Request) {
	addr := req.RemoteAddr.String()
	for _, p := range players {
		if p.addr.String() == addr {
			delete(players, p.Name)
			return
		}
	}
}

func lookupPlayer(name string, req *network.Request) (p *player, ok bool) {

	if req != nil {
		addr := req.RemoteAddr.String()
		for _, p := range players {
			if p.addr.String() == addr {
				return p, true
			}
		}
	}

	p, ok = players[name]

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

	p := newPlayer(name, req)
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

	p, ok := lookupPlayer("", req)
	// p, ok := players[req.RemoteAddr.String()]
	if !ok {
		log.Println("WARNING player is not connected", req.RemoteAddr.String(), players)
		return nil
	}

	p.deadline = resetDeadline()

	s.UpdatePlayer(p)

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
