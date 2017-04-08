package engine

import (
	"assets"
	"log"
	"phys"
	"render"

	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	window     *glfw.Window
	callbacks  []func(float32) bool
	NeedRender bool
)

//Client main function create window and initialize opengl,render engine
func Client(userfunc func()) {
	NeedRender = true

	var err error
	window, err = render.NewWindow(1200, 800, "phasarv-client")
	if err != nil {
		log.Panicln(err)
	}

	err = assets.LoadAssets("assets/textures", "assets/shaders", "assets/models")
	if err != nil {
		log.Panicln("failed load assets, reason: %s", err)
	}

	phys.Init()

	userfunc()

	LoopRender()
}

//Server mail function of network server part
func Server(userfunc func()) {
	phys.Init()

	userfunc()

	LoopServer()
}

//AddCallback add  engine.ch frame callback with deltaT as one argument
func AddCallback(funcs ...func(float32) bool) {
	for _, f := range funcs {
		callbacks = append(callbacks, f)
		// e.callbacks[len(e.callbacks)] = f
	}
}

//SetMouseCallback set function  callback each frame
func SetMouseCallback(f func(*glfw.Window, glfw.MouseButton, glfw.Action, glfw.ModifierKey)) {
	window.SetMouseButtonCallback(f)
}

//SetKeyCallback set function  callback each frame
func SetKeyCallback(f func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)) {
	window.SetKeyCallback(f)
}

//CursorPosition returned cursor position
func CursorPosition() (float32, float32) {
	x, y := window.GetCursorPos()
	return float32(x), float32(y)
}

//WindowSize returned current size of window [type: float32]
func WindowSize() (float32, float32) {
	w, h := window.GetSize()
	return float32(w), float32(h)
}
