package main

import (
	"log"
	"param"
	"runtime"
	"scene"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/renderer/forward"

	"engine"
)

func init() {
	log.SetFlags(log.Lshortfile)
	runtime.LockOSThread()
}

var (
	cursor *engine.Object
	camera *fizzle.YawPitchCamera

	light0 *forward.Light
	light1 *forward.Light
)

func main() {
	engine.NewWindow()
	engine.InitPhys(0.3)

	if err := scene.Load("scene00"); err != nil {
		log.Fatalln("failed load scene", "scene00", err)
	}

	initEnvironment()
	initCursor()
	initLocalPlayer()

	engine.MainLoop()
}

func initEnvironment() {
	light0 = engine.NewLight(100, 1, 4096)
	light1 = engine.NewLight(10, 1, 2)

	camera = engine.NewCamera(mgl32.Vec3{0, 0, 40})
	camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
	// camera.SetYawAndPitch(0, mgl32.DegToRad(-90))
}

func initCursor() {
	cursor = engine.NewPlanePoint(param.Object{
		Name: "cursor",
		Mesh: param.Mesh{Shader: "colortext2", Texture: "cursor"},
	}, mgl32.Vec3{-0.5, -0.5, 1}, mgl32.Vec3{0.5, 0.5, 1})
	cursor.Node.Core.DiffuseColor = mgl32.Vec4{0.3, 0.3, 1, 0.9}

	// engine.AddCallback(cameraMovement)
	// engine.Control = cameraMovement
}
