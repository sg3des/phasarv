package main

import (
	"client"
	"db"
	"flag"
	"log"
	"param"
	"scene"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/renderer/forward"

	"engine"
)

var (
	mode string
)

func init() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	mode = flag.Arg(0)
}

var (
	cursor *engine.Object
	camera *fizzle.YawPitchCamera
	sun    *forward.Light

	players []*Player
)

func main() {
	engine.Init()
	engine.SetKeyCallback(keyCallback)

	if mode == "network" {
		networkPlay()
	} else {
		localPlay()
	}

	scene.Load("scene00")

	initEnvironment()
	initCursor()
	initEnemies()

	engine.Loop()
}

func networkPlay() {
	client.AddRoutes(map[string]func(interface{}) interface{}{
		"CreateLocalPlayer": CreateLocalPlayer,
	})

	client.Connect("127.0.0.1:9696")
	client.Authorize("player0")
}

func localPlay() {
	CreateLocalPlayer(db.GetPlayer("player0"))
}

func initEnvironment() {
	sun = engine.NewSun()
	engine.NewLight(10, 1, 2)

	camera = engine.NewCamera(mgl32.Vec3{0, 0, 40})
	camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
}

func initCursor() {
	cursor = engine.NewPlanePoint(
		param.Object{
			Name:     "cursor",
			Mesh:     param.Mesh{Shader: "colortext2"},
			Material: param.Material{Name: "cursor", DiffColor: mgl32.Vec4{0.3, 0.3, 1, 0.9}},
		},
		mgl32.Vec3{-0.5, -0.5, 1},
		mgl32.Vec3{0.5, 0.5, 1},
	)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	if key == glfw.KeyF1 && action == glfw.Press {
		glfw.SwapInterval(1)
	}
}
