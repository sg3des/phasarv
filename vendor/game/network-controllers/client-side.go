package controllers

import (
	"engine"
	"game"
	"game/ships"
	"log"
	"materials"
	"network"
	"phys/vect"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	name        string
	c           Client
	localplayer *cliPlayer
	pshadow     *engine.Object
)

func Connect(addr string) {
	conn = network.NewConnection(c)
	err := conn.Client(addr)
	if err != nil {
		log.Fatalln(err)
	}
}

func SendAuthorize(name string) {
	if err := conn.SendMessage(Server.Authorize, Client.LoadLocalPlayer, name); err != nil {
		log.Fatalln("failed authorize", err)
	}
}

func sendLocalPlayerState(_ float32) bool {
	conn.SendMessage(
		Server.PlayerState,
		Client.LocalPlayerServerState,
		localplayer.GetClientState(),
	)
	return true
}

//
//
// Client side controllers handlers
type Client struct{}

func (Client) LoadLocalPlayer(req *network.Request) interface{} {
	s := new(ships.Ship)
	*s = req.Data.(ships.Ship)

	p := game.NewLocalPlayer(s, name)

	engine.AddCallback(sendLocalPlayerState)
	// p.Ship.Object.AddCallback(sendLocalPlayerState)

	localplayer = addCliPlayer(p)

	// localplayer = &player{p, req.RemoteAddr, time.Time{}}
	// *localplayer =
	// game.CreateLocalPlayer(localplayer)

	pshadow = &engine.Object{
		Name: "shadow",
		RI: &render.Instruction{
			MeshName:    "trapeze",
			Material:    &materials.Instruction{Name: "player", Texture: "gray", Shader: "blend", SpecLevel: 1, DiffColor: mgl32.Vec4{1, 0, 0, 0.7}},
			Shadow:      false,
			Transparent: true,
		},
	}
	pshadow.Create()

	return nil
}

func (Client) LoadPlayer(req *network.Request) interface{} {
	p, ok := req.Data.(game.Player)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	game.CreatePlayer(&p)
	addCliPlayer(&p)

	return nil
}

func (Client) RemovePlayer(req *network.Request) interface{} {
	name, ok := req.Data.(string)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	log.Println(name)

	game.RemovePlayer(name)
	delCliPlayer(name)

	return nil
}

func (Client) Enemy(req *network.Request) interface{} {
	pos := req.Data.(vect.Vect)
	log.Println("load enemy", pos)
	game.CreateEnemy()
	return nil
}

func (Client) LocalPlayerServerState(req *network.Request) interface{} {
	return nil
	s, ok := req.Data.(ServerState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	pshadow.SetPosition(s.Pos.X, s.Pos.Y)
	pshadow.SetRotation(s.Rot)

	dist := localplayer.p.Ship.Object.DistancePoint(s.Pos.X, s.Pos.Y)
	if dist > 2 {
		log.Println("WARNING! need correct position")
		localplayer.p.Ship.Object.SetPosition(s.Pos.X, s.Pos.Y)
		localplayer.p.Ship.Object.SetRotation(s.Rot)
		s.UpdatePlayer(localplayer.p)
		return nil
	}

	localplayer.p.Ship.Object.SetVelocity(s.Vel.X, s.Vel.Y)
	localplayer.p.Ship.Object.SetAngularVelocity(s.AVel)

	return nil
}

func (Client) PlayersServerState(req *network.Request) interface{} {
	states, ok := req.Data.(ServersState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	for _, s := range states {

		//skip update for localplayer
		if localplayer != nil && s.Name == localplayer.p.Name {
			// 	continue
			pshadow.SetPosition(s.Pos.X, s.Pos.Y)
			pshadow.SetRotation(s.Rot)

		}

		c, ok := lookupCliPlayer(s.Name, nil)
		// log.Println(s.Name, p, ok)
		// p, ok := game.LookupPlayer(s.Name)
		if ok {
			s.UpdatePlayer(c.p)
		} else {
			conn.SendMessage(Server.GetPlayer, Client.LoadPlayer, s.Name)
		}

		//update other players
		// for _, p := range game.Players {
		// 	if p.Name == s.Name {
		// 		p.UpdateFromServerState(s)
		// 	}
		// }
	}

	return nil
}
