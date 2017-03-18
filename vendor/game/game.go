package game

import (
	"engine"
	"log"
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

var (
	Camera *fizzle.YawPitchCamera

	//Players its clients
	Players []*Player

	//itsClient flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
	NeedRender bool
)

func CreateLocalPlayer(p *Player) {
	Camera = render.NewCamera(mgl32.Vec3{0, 0, 40})
	Camera.LookAtDirect(mgl32.Vec3{0, 0, 0})

	log.Println("CreateLocalPlayer", p.Name)
	p.Local = true
	// p := &Player{Param: paramPlayer}
	p.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.CreatePlayer()

	p.Object.SetDestroyFunc(p.Destroy)

	Players = append(Players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	p.Object.AddCallback(p.Movement, p.PlayerRotation, p.CameraMovement, p.Attack)
	engine.SetMouseCallback(p.MouseControl)
}

func CreatePlayer(p *Player) {
	log.Println("CreatePlayer", p.Name)

	// p := &Player{Param: paramPlayer}
	p.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.CreatePlayer()

	p.Object.SetDestroyFunc(p.Destroy)

	Players = append(Players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	p.Object.AddCallback(p.Movement, p.PlayerRotation, p.Attack)
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
			p.Object.Remove()
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
