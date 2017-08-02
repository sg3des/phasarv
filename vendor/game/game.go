package game

import (
	"engine"
	"game/ships"
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

var (
	Camera *fizzle.YawPitchCamera

	//Players contains only human playble ships
	Players []*Player

	//NeedRender flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
	NeedRender bool
)

type Player struct {
	Name string
	Ship *ships.Ship

	Target *Player

	targetAngle  float32
	respawnPoint mgl32.Vec2
}

func CreateLocalPlayer(p *Player) {
	Camera = render.NewCamera(mgl32.Vec3{0, 0, 40})
	Camera.LookAtDirect(mgl32.Vec3{0, 0, 0})

	p.Ship.Local = true
	// p := &Player{Param: paramPlayer}
	p.Ship.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.Ship.Create()

	p.Ship.Object.SetDestroyFunc(p.Ship.Destroy)

	Players = append(Players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	p.Ship.Object.AddCallback(p.Ship.Movement, p.Ship.Rotation, p.CameraMovement, p.Ship.Attack)
	engine.SetMouseCallback(p.Ship.MouseControl)
}

func NewLocalPlayer(s *ships.Ship, name string) *Player {
	p := &Player{
		Name: name,
		Ship: s,
	}

	CreateLocalPlayer(p)

	return p
}

func NewPlayer(name string, s *ships.Ship, a mgl32.Vec2) *Player {
	p := &Player{
		Name: name,
		Ship: s,
	}

	CreatePlayer(p)

	return p
}

func CreatePlayer(p *Player) {
	p.Ship.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.Ship.Create()

	p.Ship.Object.SetDestroyFunc(p.Ship.Destroy)

	Players = append(Players, p)

	p.Ship.Object.AddCallback(p.Ship.Movement, p.Ship.Rotation, p.Ship.Attack)

	return
}

func LookupPlayer(name string) (*Player, bool) {
	for _, p := range Players {
		if p.Name == name {
			return p, true
		}
	}

	return nil, false
}

func RemovePlayer(name string) {
	for i, p := range Players {
		if p.Name == name {
			p.Ship.Object.Remove()
			RemovePlayerN(i)
			return
		}
	}
}

func RemovePlayerN(i int) {
	Players[i] = Players[len(Players)-1]
	Players[len(Players)-1] = nil
	Players = Players[:len(Players)-1]
}
