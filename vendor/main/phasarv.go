package main

import (
	"log"
	"param"
	"scene"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"

	"engine"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	cursor *engine.Object
	camera *fizzle.YawPitchCamera

	players []*Player
)

func main() {
	if err := scene.Load("scene00"); err != nil {
		log.Fatalln("failed load scene", "scene00", err)
	}

	initEnvironment()
	initCursor()
	initLocalPlayer()

	initEnemies()

	engine.Loop()
}

func initEnvironment() {
	engine.NewSun()
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
