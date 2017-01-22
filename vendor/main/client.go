package main

import (
	"encoding/gob"
	"game"
	"log"
	"network"
	"phys/vect"
	"time"
)

var (
	p = &game.Player{}
	c *network.Connection
)

func Connect(addr string) {
	gob.Register(game.Player{})
	// gob.Register(game.Weapon{})
	gob.Register(vect.Vect{})
	gob.Register(game.NetPacket{})
	// gob.Register(point.P{})
	// gob.Register(point.Param{})
	// gob.Register(render.Renderable{})
	// gob.Register(render.Instruction{})
	// gob.Register(phys.Instruction{})
	// gob.Register(engine.Object{})
	// gob.Register(mgl32.Vec3{})
	// gob.Register(materials.Instruction{})

	c = network.NewHandlers(map[string]network.Handler{
		"loadLocalPlayer":   loadLocalPlayer,
		"loadEnemy":         loadEnemy,
		"getPlayerPosition": getPlayerPosition,
	})
	if err := c.Client(addr); err != nil {
		log.Fatalln(err)
	}

}

func Authorize(name string) {
	if err := c.SendMessage("auth", "loadLocalPlayer", name); err != nil {
		log.Fatalln("failed authorize", err)
	}
}

func loadLocalPlayer(req *network.Request) interface{} {
	*p = req.Data.(game.Player)
	game.CreateLocalPlayer(p)
	p.Object.AddCallback(sendCursorPosition)
	return nil
}

var sendCursorPositionTime time.Time

func sendCursorPosition(dt float32) {
	t := time.Now()
	if t.After(sendCursorPositionTime) {
		c.SendMessage("playersCursor", "getPlayerPosition", p.Cursor.PositionVect())
		sendCursorPositionTime = t.Add(time.Millisecond * 100)
	}
}

func getPlayerPosition(req *network.Request) interface{} {
	np, ok := req.Data.(game.NetPacket)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	// log.Println(p.Object.PositionVect(), p.Cursor.PositionVect())
	dist := p.Object.DistancePoint(np.Pos.X, np.Pos.Y)
	if dist > 1 {
		log.Println("correct position")
		p.Object.SetPosition(np.Pos.X, np.Pos.Y)
	}
	// log.Println("update position")
	// p.Object.SetPosition(v.X, v.Y)
	// pos := p.Object.PositionVect()
	// pos.Sub(np.Pos)
	// np.Vel.Add(pos)
	// pos.Add(np.Vel)

	p.Object.SetVelocity(np.Vel.X, np.Vel.Y)
	p.Object.SetAngularVelocity(np.AVel)
	log.Println(dist, p.Object.PositionVect(), np.Pos)
	// }

	return nil
}

func loadEnemy(req *network.Request) interface{} {
	pos := req.Data.(vect.Vect)
	log.Println("load enemy", pos)
	game.CreateEnemy()
	return nil
}
