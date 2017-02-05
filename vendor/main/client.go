package main

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
	localplayer = &game.Player{}
	pshadow     *engine.Object
	conn        *network.Connection
)

func Connect(addr string) {
	game.RegisterNetworkTypes()

	conn = network.NewHandlers(map[string]network.Handler{
		"loadLocalPlayer": loadLocalPlayer,
		"loadPlayer":      loadPlayer,
		"loadEnemy":       loadEnemy,
		"getServerState":  getServerState,
		"getServersState": getServersState,
	})
	if err := conn.Client(addr); err != nil {
		log.Fatalln(err)
	}

}

func Authorize(name string) {
	if err := conn.SendMessage("auth", "loadLocalPlayer", name); err != nil {
		log.Fatalln("failed authorize", err)
	}
}

func loadLocalPlayer(req *network.Request) interface{} {
	*localplayer = req.Data.(game.Player)
	game.CreateLocalPlayer(localplayer)
	localplayer.Object.AddCallback(sendClientState)

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

func loadPlayer(req *network.Request) interface{} {
	p, ok := req.Data.(game.Player)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	game.CreatePlayer(&p)

	return nil
}

// var sendCursorPositionTime time.Time

func sendClientState(dt float32) {
	// t := time.Now()
	// if t.After(sendCursorPositionTime) {
	// state := game.ClientState{
	// 	CurPos: p.Cursor.PositionVect(),
	// 	LW:     p.LeftWeapon.ToShoot,
	// 	RW:     p.RightWeapon.ToShoot,
	// }

	conn.SendMessage("clientState", "getServerState", localplayer.GetClientState())
	// sendCursorPositionTime = t.Add(time.Millisecond * 50)
	// }
}

func getServerState(req *network.Request) interface{} {
	s, ok := req.Data.(game.ServerState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	pshadow.SetPosition(s.Pos.X, s.Pos.Y)
	pshadow.SetRotation(s.Rot)

	dist := localplayer.Object.DistancePoint(s.Pos.X, s.Pos.Y)
	if dist > 2 {
		log.Println("WARNING! need correct position")
		localplayer.UpdateFromServerState(s)
		return nil
	}

	localplayer.Object.SetVelocity(s.Vel.X, s.Vel.Y)
	localplayer.Object.SetAngularVelocity(s.AVel)

	return nil
}

func getServersState(req *network.Request) interface{} {
	states, ok := req.Data.(game.ServersState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
	}

	for _, s := range states {

		//skip update for localplayer
		if s.Name == localplayer.Name {
			continue
		}

		p, ok := game.LookupPlayer(s.Name)
		if ok {
			p.UpdateFromServerState(s)
		} else {
			conn.SendMessage("getPlayer", "loadPlayer", s.Name)
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

func loadEnemy(req *network.Request) interface{} {
	pos := req.Data.(vect.Vect)
	log.Println("load enemy", pos)
	game.CreateEnemy()
	return nil
}
