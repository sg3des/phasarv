package controllers

//ServerState structure of standard network packet
import (
	"game/players"
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

// func init() {
// 	gob.Register(players.User{})

// 	gob.Register(ClientState{})
// 	gob.Register(ServerState{})
// 	gob.Register(ServersState{})
// }

type user struct {
	u *players.User
	p *players.Player

	addr *net.UDPAddr

	deadline time.Time
}

func (c *user) GetClientState() ClientState {
	return ClientState{
		CurPos: c.p.Ship.Cursor.PositionVect(),
		LW:     c.p.Ship.LeftWeapon.ToShoot,
		RW:     c.p.Ship.RightWeapon.ToShoot,
	}
}

func (c *user) GetServerState() ServerState {
	return ServerState{
		Name: c.p.Name,

		Vel:  c.p.Ship.Object.Velocity(),
		AVel: c.p.Ship.Object.AngularVelocity(),

		Pos: c.p.Ship.Object.PositionVect(),
		Rot: c.p.Ship.Object.Rotation(),

		HP: c.p.Ship.CurrParam.Health,

		ClientState: c.GetClientState(),
	}
}

func (s ClientState) UpdatePlayer(p *players.Player) {
	p.Ship.Cursor.SetPosition(s.CurPos.X, s.CurPos.Y)
	p.Ship.CursorOffset = p.Ship.Cursor.PositionVect()
	p.Ship.CursorOffset.Sub(p.Ship.Object.PositionVect())

	p.Ship.LeftWeapon.ToShoot = s.LW
	p.Ship.RightWeapon.ToShoot = s.RW
}

func (s ServerState) UpdatePlayer(p *players.Player) {
	p.Ship.Object.SetPosition(s.Pos.X, s.Pos.Y)
	p.Ship.Object.SetRotation(s.Rot)
	p.Ship.Object.SetVelocity(s.Vel.X, s.Vel.Y)
	p.Ship.Object.SetAngularVelocity(s.AVel)
	p.Ship.CurrParam.Health = s.HP
	s.ClientState.UpdatePlayer(p)
	// p.UpdateFromClientState(s.ClientState)
}
