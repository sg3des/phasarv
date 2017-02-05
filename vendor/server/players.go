package main

import (
	"db"
	"game"
	"log"
	"network"
)

var players map[string]*game.Player

func auth(req *network.Request) interface{} {
	name, ok := req.Data.(string)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
	}
	log.Println("auth", name)

	p := db.GetPlayer(name)
	p.CreatePlayer()
	p.Object.AddCallback(p.ClientCursor, p.Movement, p.PlayerRotation)
	players[req.RemoteAddr.String()] = p

	return db.GetPlayer(name)
}

func getPlayer(req *network.Request) interface{} {
	name, ok := req.Data.(string)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
	}

	return db.GetPlayer(name)
}

func clientState(req *network.Request) interface{} {
	s, ok := req.Data.(game.ClientState)
	if !ok {
		log.Println("WARNING! recieve data is not correct")
		return nil
	}

	p, ok := players[req.RemoteAddr.String()]
	if !ok {
		log.Println("WARNING player is not connected", req.RemoteAddr.String(), players)
		return nil
	}

	p.UpdateFromClientState(s)

	// log.Println(p.Object.PositionVect(), p.Cursor.PositionVect())

	return game.ServerState{
		Vel:  p.Object.Velocity(),
		AVel: p.Object.AngularVelocity(),
		Pos:  p.Object.PositionVect(),
		Rot:  p.Object.Rotation(),
	}
}

func sendServersState(float32) bool {
	if len(players) == 0 {
		return true
	}

	var states game.ServersState
	for _, p := range players {
		states = append(states, p.GetServerState())
	}
	s.Broadcast("getServersState", "", states)

	return true
}
