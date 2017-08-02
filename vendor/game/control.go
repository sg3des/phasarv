package game

import (
	"engine"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

func (p *Player) CameraMovement(dt float32) {
	render.SetCameraPosition(p.Ship.Object.Position())

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()
	campos := render.GetCameraPosition()

	d := h + campos.Z()

	x = (x-w/2)/d*campos.Z() + campos.X()
	y = (h/2-y)/d*campos.Z() + campos.Y()

	// p.Ship.Cursor.SetPosition(getCursorPos(x, y, w, h, render.GetCameraPosition()))
	p.Ship.Cursor.SetPosition(x, y)
	p.Ship.LeftWeapon.CursorPos = mgl32.Vec2{x, y}
	p.Ship.RightWeapon.CursorPos = mgl32.Vec2{x, y}
}

func getCursorPos(x, y, w, h float32, campos mgl32.Vec3) (float32, float32) {
	d := h + campos.Z()

	x = (x-w/2)/d*campos.Z() + campos.X()
	y = (h/2-y)/d*campos.Z() + campos.Y()

	return x, y
}
