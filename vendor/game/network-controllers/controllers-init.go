package controllers

//ServerState structure of standard network packet
import (
	"encoding/gob"
	"game"
	"net"
	"network"
	"phys/vect"
	"time"
)

var (
	conn             *network.Connection
	durationDeadline = time.Duration(1e9)
)

type ServersState []ServerState
type ServerState struct {
	Name string

	Vel  vect.Vect
	AVel float32

	Pos vect.Vect
	Rot float32

	HP float32

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

type cliPlayer struct {
	p *game.Player

	addr *net.UDPAddr

	deadline time.Time
}

func (c *cliPlayer) GetClientState() ClientState {
	return ClientState{
		CurPos: c.p.Cursor.PositionVect(),
		LW:     c.p.LeftWeapon.ToShoot,
		RW:     c.p.RightWeapon.ToShoot,
	}
}

func (c *cliPlayer) GetServerState() ServerState {
	return ServerState{
		Name: c.p.Name,

		Vel:  c.p.Object.Velocity(),
		AVel: c.p.Object.AngularVelocity(),

		Pos: c.p.Object.PositionVect(),
		Rot: c.p.Object.Rotation(),

		HP: c.p.CurrParam.Health,

		ClientState: c.GetClientState(),
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
	p.CurrParam.Health = s.HP
	s.ClientState.UpdatePlayer(p)
	// p.UpdateFromClientState(s.ClientState)
}
