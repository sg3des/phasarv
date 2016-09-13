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
	engine.InitPhys(0.5)

	if err := scene.Load("scene00"); err != nil {
		log.Fatalln("failed load scene", "scene00", err)
	}
	initEnvironment()
	initCursor()
	newPlayer()

	engine.MainLoop()
}

func initEnvironment() {
	light0 = engine.NewLight(100, 1, 4096)
	light1 = engine.NewLight(10, 1, 2)

	camera = engine.NewCamera(mgl32.Vec3{0, 0, 40})
	camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
	// camera.SetYawAndPitch(0, mgl32.DegToRad(-90))
}

// func loadScene() {

// 	engine.InitPhys(0.5)

// 	ground := engine.NewObject(param.Object{
// 		Name: "ground",
// 		Mesh: param.Mesh{Model: "plane", Texture: "grass", Shader: "diffuse_texbumped_shadows"},
// 		Pos:  param.Pos{0, 0, -10},
// 	}, nil)
// 	ground.Node.Core.SpecularColor = mgl32.Vec4{0, 0, 0, 0}

// 	engine.CreateCurve(fizzle.X | fizzle.Y).Location = mgl32.Vec3{5, 13, 1}
// 	// engine.CreateCurve(fizzle.X | fizzle.Z).Location = mgl32.Vec3{0, 13, 1}
// 	// engine.CreateCurve(fizzle.Y | fizzle.Z).Location = mgl32.Vec3{-5, 13, 1}
// 	// engine.CreateCurve(fizzle.X | fizzle.Y | fizzle.Z).Location = mgl32.Vec3{-10, 0, 1}
// 	// return

// 	hb := []param.Art{newHealthBar(100)}

// 	// tree := engine.NewObject(param.Object{
// 	// 	Name:        "box0",
// 	// 	Mesh:        param.Mesh{Model: "tree", Texture: "gray", Shader: "colortext2"},
// 	// 	Pos:         param.Pos{10, 5, 0},
// 	// 	PH:          param.Phys{H: 2, W: 2, Mass: 10, Group: 1},
// 	// 	Transparent: true,
// 	// }, hb)
// 	// tree.Node.Core.DiffuseColor = mgl32.Vec4{0.1, 0.5, 0.1, 0.9}
// 	// tree.Node.Scale = mgl32.Vec3{0.04, 0.04, 0.04}

// 	engine.NewObject(param.Object{
// 		Name: "box1",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{10, -8, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 1, Group: 1},
// 	}, hb).SetRotation(1.5)

// 	engine.NewObject(param.Object{
// 		Name: "box1",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{10, -3, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 1000, Group: 1},
// 	}, hb)

// 	engine.NewObject(param.Object{
// 		Name: "box0",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{-10, 5, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 100, Group: 1},
// 	}, hb)

// 	engine.NewObject(param.Object{
// 		Name: "box1",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{-10, -8, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 100, Group: 1},
// 	}, hb)

// 	engine.NewObject(param.Object{
// 		Name: "box1",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{1, -10, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 100, Group: 1},
// 	}, hb)

// 	engine.NewObject(param.Object{
// 		Name: "box1",
// 		Mesh: param.Mesh{Model: "t", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
// 		Pos:  param.Pos{1, 10, 0},
// 		PH:   param.Phys{H: 2, W: 2, Mass: 100, Group: 1},
// 	}, hb)
// }

func initCursor() {
	cursor = engine.NewPlane(param.Object{
		Name: "cursor",
		Mesh: param.Mesh{Shader: "colortext2", Texture: "cursor"},
	}, 1, 1)
	cursor.Node.Core.DiffuseColor = mgl32.Vec4{0.3, 0.3, 1, 0.9}

	// cursor.Callback = cursorCallback

	engine.Control = cameraMovement

}

// func cursorCallback(c *engine.Object, dt float32) {
// 	c.Childs
// }

func cameraMovement() {
	pp := localPlayer.Node.Location

	cp := camera.GetPosition()
	camera.SetPosition(pp.X(), pp.Y(), cp.Z())

	light0.Position = mgl32.Vec3{pp.X() - 5, pp.Y() + 5, 30}
	// light0.Position = mgl32.Vec3{pp.X() - 10, pp.Y() + 10, 30}
	// light1.Position = mgl32.Vec3{pp.X() - 3, pp.Y() + 3, -3}

	xfloat, yfloat := engine.Window.GetCursorPos()
	width, height := engine.Window.GetSize()
	// cameraPos := camera.GetPosition()

	x, y := getCursorPos(float32(xfloat), float32(yfloat), width, height, cp)

	cursor.Node.Location = mgl32.Vec3{x, y, 0}
	// log.Println(x, y, xfloat, yfloat)
}
