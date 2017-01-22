package main

import (
	"game"
	"log"
	"network"
	"phys/vect"
)

var players map[string]*game.Player

func playersCursor(req *network.Request) interface{} {
	v, ok := req.Data.(vect.Vect)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	p, ok := players[req.RemoteAddr.String()]
	if !ok {
		log.Println("WARNING player is not connected", req.RemoteAddr.String(), players)
		return nil
	}

	p.Cursor.SetPosition(v.X, v.Y)

	// log.Println(p.Object.PositionVect(), p.Cursor.PositionVect())

	return game.NetPacket{Vel: p.Object.Velocity(), AVel: p.Object.AngularVelocity(), Pos: p.Object.PositionVect()}
}
