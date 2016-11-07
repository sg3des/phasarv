package main

import (
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
	engine.Init(local)
	// engine.SetKeyCallback(keyCallback)
	// initEnvironment()

	// scene.Load("scene00")

	// // // time.Sleep(300 * time.Millisecond)

	// if mode == "client" {
	// 	networkPlay()
	// } else {
	// 	localPlay()
	// }

	// initCursor()

	// time.Sleep(10 * time.Second)

	// engine.Loop()
}

func local() {
	engine.SetKeyCallback(keyCallback)
	initEnvironment()
	initCursor()

	scene.Load("scene00")

	// // time.Sleep(300 * time.Millisecond)

	if mode == "client" {
		networkPlay()
	} else {
		localPlay()
	}

}

func networkPlay() {
	// client.AddRoutes(map[string]func(interface{}) interface{}{
	// 	"CreateLocalPlayer": CreateLocalPlayer,
	// })

	Connect("127.0.0.1:9696")
	Authorize("player0")
}

func localPlay() {
	// createEnemy(10, 15)
	CreateLocalPlayer(db.GetPlayer("player0"))
}

func initEnvironment() {
	engine.NewLight(engine.ParamLight{
		Sun:        true,
		Pos:        mgl32.Vec3{-30, 30, 60},
		Strength:   0.7,
		Specular:   0.5,
		ShadowSize: 8192,
	})
	// sun = engine.NewSun()

	engine.NewLight(engine.ParamLight{
		Pos:        mgl32.Vec3{1, 1, -2},
		Strength:   0.1,
		Specular:   0,
		ShadowSize: 2,
	})

	camera = engine.NewCamera(mgl32.Vec3{0, 0, 40})
	camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
}

func initCursor() {
	cursor = engine.NewObject(
		param.Object{
			Name:     "cursor",
			Mesh:     param.Mesh{Model: "plane", X: 1, Y: 1},
			Material: param.Material{Name: "cursor", Shader: "colortext2", DiffColor: mgl32.Vec4{0.3, 0.3, 1, 0.9}},
		},

		// mgl32.Vec3{-0.5, -0.5, 1},
		// mgl32.Vec3{0.5, 0.5, 1},
	)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	if key == glfw.KeyF1 && action == glfw.Press {
		glfw.SwapInterval(1)
	}
	if key == glfw.KeySpace && action == glfw.Press {
		for _, p := range players {
			p.Destroy()
		}
	}
}
