package main

import (
	"db"
	"flag"
	"game"
	"log"
	"math/rand"
	"render"
	"scene"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/renderer/forward"

	"engine"

	controllers "network-controllers"
)

var (
	client bool
	cursor *engine.Object
	camera *fizzle.YawPitchCamera
	sun    *forward.Light
)

func init() {
	// runtime.LockOSThread()
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	if flag.Arg(0) == "client" {
		client = true
	}
	// mode =
}

func main() {
	game.NeedRender = true
	engine.Client(local)

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

	scene.Load("scene00")

	if client {
		go networkPlay()
	} else {
		localPlay()
	}
}

func networkPlay() {
	controllers.Connect("127.0.0.1:9696")
	controllers.SendAuthorize(randomName())
}

func localPlay() {
	game.CreateLocalPlayer(db.GetPlayer(randomName()))
	// for i := 0; i < 10; i++ {
	// 	game.CreateEnemy()
	// }
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomName() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
