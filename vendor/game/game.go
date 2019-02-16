package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	enter  *Enter
	hangar *Hangar
	battle *Battle

	//NeedRender flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
	NeedRender bool

	SinglePlay bool
)

func Start() {
	enter = NewEnter()
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}
