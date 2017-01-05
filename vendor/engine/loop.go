package engine

import (
	"engine/frames"
	"fmt"
	"phys"
	"render"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/eweygewey"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

func ReadMemory() {
	runtime.ReadMemStats(&mem)
	fmt.Printf("%dmb %dmb %dmb %dmb\n", mem.Alloc/1048576, mem.TotalAlloc/1048576, mem.HeapAlloc/1048576, mem.HeapSys/1048576)
}

//Loop is function where:  poll events, calculate physics, rendered objects and effects. This is infinite loop, and should be run late from user space
func Loop() {
	// time.Sleep(500 * time.Millisecond)
	var dt, Rdt float32
	var fps, Rfps int

	wnd := e.ui.NewWindow("phys", 0.79, 0.99, 0.2, 0, func(wng *eweygewey.Window) {
		wng.Text("phys")
		wng.StartRow()
		wng.Text(fmt.Sprintf("fps: %d dt: %f", fps, dt))

	})
	wnd.ShowTitleBar = false
	wnd.IsMoveable = false
	wnd.AutoAdjustHeight = true

	Rwnd := e.ui.NewWindow("render", 0.79, 0.91, 0.2, 0, func(wng *eweygewey.Window) {
		wng.Text("render")
		wng.StartRow()
		wng.Text(fmt.Sprintf("fps: %d dt: %f", Rfps, Rdt))
	})
	Rwnd.ShowTitleBar = false
	Rwnd.IsMoveable = false
	Rwnd.AutoAdjustHeight = true

	PhysFrame := frames.NewFrame()
	Rendframe := frames.NewFrame()

	// loopCreator()

	go func() {
		for {
			dt, fps = PhysFrame.Next()

			loopCallbacks(dt)
			phys.NextFrame(dt)
		}
	}()

	for !e.window.ShouldClose() {
		Rdt, Rfps = Rendframe.Next()
		loopPhys()
		// loopCreator()

		// loopRenderShadows()
		w, h := e.window.GetFramebufferSize()

		// clear the screen and reset our viewport
		e.gfx.Viewport(0, 0, int32(w), int32(h))
		// e.gfx.ClearColor(0.5, 0.5, 0.5, 1)
		e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		render.NextFrame(mgl32.DegToRad(50), float32(w)/float32(h))

		e.ui.Construct(float64(dt))
		e.ui.Draw()

		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loopPhys() {
	// log.Println(dt)
	// phys.space.Step(dt)

	// log.Println(len(Objects))
	for o := range Objects {

		if o.renderable.Body == nil {
			// log.Println("skip:", o.Name)
			continue
		}

		// log.Println("rend:", o.Name)

		// update position
		o.renderable.Body.Location = o.PositionVec3()

		// log.Println(o.Name, o.PositionVec3())

		// update rotation
		ang := o.Rotation()
		if o.RollAngle != 0 {
			q := mgl32.AnglesToQuat(0, 0, ang, 1).Mul(mgl32.AnglesToQuat(o.RollAngle, 0, 0, 1))
			o.renderable.Body.LocalRotation = q

			shape := o.Shape.GetAsBox()
			shape.Width = o.PI.W - o.PI.W*o.ShapeWidthPercent()
			shape.UpdatePoly()
		} else {
			o.renderable.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, 1)
		}
	}
}

func loopCallbacks(dt float32) {
	for i, f := range e.callbacks {
		if f != nil && !f(dt) {
			delete(e.callbacks, i)
		}
	}

	for o := range Objects {
		for _, f := range o.Callbacks {
			f(dt)
		}
	}
}

// func loopCreator() {
// 	for l := range lights {
// 		if l.LightNode == nil {
// 			log.Println("create light, direct:", l.Direct)
// 			l.create()
// 		}
// 	}

// 	for _, o := range Scene {
// 		if o.Body == nil {
// 			log.Println("create scene object: ", o.Name)
// 			o.createRenderable(o.Param)
// 		}
// 	}

// 	for o := range Objects {
// 		if o.Body == nil {
// 			log.Println("create object: ", o.Name)
// 			o.createRenderable(o.Param)

// 			// for _, art := range o.Arts {
// 			// 	// if art.Art == nil {
// 			// 	art.createNode(art.Param)
// 			// 	// }
// 			// }
// 			// for _, art := range o.ArtRotate {
// 			// 	if art.Body == nil {
// 			// 		log.Println("create art:", art.Name, art.Node)
// 			// 		art.createNode(art.Param)
// 			// 	}
// 			// }

// 			// for _, art := range o.ArtStatic {
// 			// 	if art.Body == nil {
// 			// 		log.Println("create art:", art.Name, art.Node)
// 			// 		art.createNode(art.Param)
// 			// 	}
// 			// }
// 		}
// 	}
// }

// func loopUI(dt float32) {
// 	// e.gfx.Disable(graphicsprovider.DEPTH_TEST)
// 	// e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

// 	e.ui.Construct(float64(dt))
// 	e.ui.Draw()

// 	// e.gfx.Disable(graphicsprovider.SCISSOR_TEST)
// 	// e.gfx.Enable(graphicsprovider.DEPTH_TEST)
// }

// func loopRender() {
// 	w, h := e.window.GetFramebufferSize()

// 	// clear the screen and reset our viewport
// 	e.gfx.Viewport(0, 0, int32(w), int32(h))
// 	// e.gfx.ClearColor(0.5, 0.5, 0.5, 1)
// 	e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

// 	render.RenderFrame(mgl32.DegToRad(50), float32(w)/float32(h))

// 	// // make the projection and view matrixes
// 	// perspective := mgl32.Perspective(mgl32.DegToRad(50), float32(w)/float32(h), 1.0, 100.0)
// 	// view := e.camera.GetViewMatrix()

// 	// // render object
// 	// for o, b := range Objects {
// 	// 	if b && o.Node != nil && !o.Transparent {
// 	// 		// log.Println("render", o.Name)
// 	// 		o.Render(perspective, view, e.camera)
// 	// 	}
// 	// }

// 	// // render scene
// 	// for _, o := range Scene {
// 	// 	if o.Node != nil && !o.Transparent {
// 	// 		o.Render(perspective, view, e.camera)
// 	// 	}
// 	// }

// 	// // render transparent objects
// 	// for o, b := range Objects {
// 	// 	if b && o.Node != nil && o.Transparent {
// 	// 		o.Render(perspective, view, e.camera)
// 	// 	}
// 	// }

// 	// // render transparent scene objects
// 	// for _, o := range Scene {
// 	// 	if o.Node != nil && o.Transparent {
// 	// 		o.Render(perspective, view, e.camera)
// 	// 	}
// 	// }

// }

// func loopRenderShadows() {
// 	e.render.StartShadowMapping()
// 	lightCount := e.render.GetActiveLightCount()
// 	for i := 0; i < lightCount; i++ {
// 		// get lights with shadow maps
// 		lightToCast := e.render.ActiveLights[i]
// 		if lightToCast.ShadowMap == nil {
// 			continue
// 		}

// 		// engine.able the light to cast shadows
// 		e.render.EnableShadowMappingLight(lightToCast)
// 		for o, b := range Objects {
// 			if b && o.Node != nil && o.Shadow {
// 				e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
// 			}
// 		}

// 		for _, o := range Scene {
// 			e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
// 		}
// 	}
// 	e.render.EndShadowMapping()
// }

//
