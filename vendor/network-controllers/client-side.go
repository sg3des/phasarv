package controllers

import (
	"engine"
	"game"
	"log"
	"materials"
	"network"
	"phys/vect"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	c           Client
	localplayer = &game.Player{}
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

func sendLocalPlayerState(dt float32) {
	conn.SendMessage(
		Server.PlayerState,
		Client.LocalPlayerServerState,
		GetClientState(localplayer),
	)
}

//
//
// Client side controllers handlers
type Client struct{}

func (Client) LoadLocalPlayer(req *network.Request) interface{} {
	*localplayer = req.Data.(game.Player)
	game.CreateLocalPlayer(localplayer)
	localplayer.Object.AddCallback(sendLocalPlayerState)

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

	return nil
}

func (Client) Enemy(req *network.Request) interface{} {
	pos := req.Data.(vect.Vect)
	log.Println("load enemy", pos)
	game.CreateEnemy()
	return nil
}

func (Client) LocalPlayerServerState(req *network.Request) interface{} {
	s, ok := req.Data.(ServerState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	pshadow.SetPosition(s.Pos.X, s.Pos.Y)
	pshadow.SetRotation(s.Rot)

	dist := localplayer.Object.DistancePoint(s.Pos.X, s.Pos.Y)
	if dist > 2 {
		log.Println("WARNING! need correct position")
		s.UpdatePlayer(localplayer)
		return nil
	}

	localplayer.Object.SetVelocity(s.Vel.X, s.Vel.Y)
	localplayer.Object.SetAngularVelocity(s.AVel)

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
		if s.Name == localplayer.Name {
			continue
		}

		p, ok := game.LookupPlayer(s.Name)
		if ok {
			s.UpdatePlayer(p)
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
