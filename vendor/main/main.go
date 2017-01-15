package main

import (
	"db"
	"flag"
	"game"
	"log"
	"render"
	"runtime"
	"scene"
	"time"

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
	runtime.LockOSThread()
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	mode = flag.Arg(0)
}

var (
	cursor *engine.Object
	camera *fizzle.YawPitchCamera
	sun    *forward.Light
)

func main() {
	engine.Init(local)
	// local()
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

	time.Sleep(time.Second)
	scene.Load("scene00")

	go localPlay()

	// // time.Sleep(300 * time.Millisecond)

	// if mode == "client" {
	// 	networkPlay()
	// } else {
	// 	localPlay()
	// }

}

func networkPlay() {
	// client.AddRoutes(map[string]func(interface{}) interface{}{
	// 	"CreateLocalPlayer": CreateLocalPlayer,
	// })

	Connect("127.0.0.1:9696")
	Authorize("player0")
}

func localPlay() {
	game.CreateLocalPlayer(db.GetPlayer("player0"))
	for i := 0; i < 10; i++ {
		game.CreateEnemy()
		log.Println(i)
	}
	// log.Println("created")

}

func initEnvironment() {
	(&render.Light{
		Direct:     true,
		Pos:        mgl32.Vec3{-30, 30, 100},
		Dir:        mgl32.Vec3{30, -30, -100},
		Strength:   0.7,
		Specular:   0.5,
		ShadowSize: 8192,
	}).Create()
	// engine.NewSun()

	(&render.Light{
		Pos:        mgl32.Vec3{1, 1, 2},
		Strength:   0.1,
		Specular:   0.1,
		ShadowSize: 2,
	}).Create()

	camera = engine.NewCamera(mgl32.Vec3{0, 0, 40})
	camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	if key == glfw.KeyF1 && action == glfw.Press {
		glfw.SwapInterval(1)
	}
	if key == glfw.KeySpace && action == glfw.Press {
		for _, p := range game.Players {
			p.Destroy()
		}
	}
}
