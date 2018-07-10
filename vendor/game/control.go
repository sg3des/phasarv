package game

import (
	"engine"
	"render"
)

func (p *Player) CameraMovement(dt float32) {
	render.SetCameraPosition(p.Ship.Object.Position())

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()
	campos := render.GetCameraPosition()

	d := h + campos.Z()

	x = (x-w/2)/d*campos.Z() + campos.X()
	y = (h/2-y)/d*campos.Z() + campos.Y()

	p.Ship.Cursor.SetPosition(x, y)
	if p.Ship.LeftWeapon != nil {
		p.Ship.LeftWeapon.UpdateCursor(x, y)
	}
	if p.Ship.RightWeapon != nil {
		p.Ship.RightWeapon.UpdateCursor(x, y)
	}
}
