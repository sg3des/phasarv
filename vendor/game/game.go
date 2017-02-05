package game

import (
	"encoding/gob"
	"phys/vect"
)

var (
	//Players its clients
	Players []*Player

	//Render flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
	Render bool
)

func LookupPlayer(name string) (*Player, bool) {
	for _, p := range Players {
		if p.Name == name {
			return p, true
		}
	}

	return nil, false
}

//ServerState structure of standard network packet
type ServerState struct {
	Name string

	Vel  vect.Vect
	AVel float32

	Pos vect.Vect
	Rot float32

	ClientState
}

type ServersState []ServerState

type ClientState struct {
	CurPos vect.Vect
	LW     bool
	RW     bool
}

func RegisterNetworkTypes() {
	gob.Register(Player{})

	gob.Register(ServerState{})
	gob.Register(ClientState{})
	gob.Register(ServersState{})
}
