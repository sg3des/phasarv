package controllers

//ServerState structure of standard network packet
import (
	"encoding/gob"
	"game"
	"network"
	"phys/vect"
)

var conn *network.Connection

type ServersState []ServerState
type ServerState struct {
	Name string

	Vel  vect.Vect
	AVel float32

	Pos vect.Vect
	Rot float32

	ClientState
}

type ClientState struct {
	CurPos vect.Vect
	LW     bool
	RW     bool
}

func init() {
	gob.Register(game.Player{})

	gob.Register(ClientState{})
	gob.Register(ServerState{})
	gob.Register(ServersState{})
}

func GetClientState(p *game.Player) ClientState {
	return ClientState{
		CurPos: p.Cursor.PositionVect(),
		LW:     p.LeftWeapon.ToShoot,
		RW:     p.RightWeapon.ToShoot,
	}
}

func GetServerState(p *game.Player) ServerState {
	return ServerState{
		Name: p.Name,

		Vel:  p.Object.Velocity(),
		AVel: p.Object.AngularVelocity(),

		Pos: p.Object.PositionVect(),
		Rot: p.Object.Rotation(),

		ClientState: GetClientState(p),
	}
}

func (s ClientState) UpdatePlayer(p *game.Player) {
	p.Cursor.SetPosition(s.CurPos.X, s.CurPos.Y)
	p.CursorOffset = p.Cursor.PositionVect()
	p.CursorOffset.Sub(p.Object.PositionVect())

	p.LeftWeapon.ToShoot = s.LW
	p.RightWeapon.ToShoot = s.RW
}

func (s ServerState) UpdatePlayer(p *game.Player) {
	p.Object.SetPosition(s.Pos.X, s.Pos.Y)
	p.Object.SetRotation(s.Rot)
	p.Object.SetVelocity(s.Vel.X, s.Vel.Y)
	p.Object.SetAngularVelocity(s.AVel)
	s.ClientState.UpdatePlayer(p)
	// p.UpdateFromClientState(s.ClientState)
}
