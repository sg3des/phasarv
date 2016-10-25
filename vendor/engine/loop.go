package engine

import (
	"engine/frames"
	"fmt"
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
	var dt float32
	var fps int

	wnd := e.ui.NewWindow("test", 0.79, 0.99, 0.2, 0, func(wng *eweygewey.Window) {
		wng.Text(fmt.Sprintf("fps: %d dt: %f", fps, dt))
	})
	wnd.ShowTitleBar = false
	wnd.IsMoveable = false
	wnd.AutoAdjustHeight = true

	frame := frames.NewFrame()
	for !e.window.ShouldClose() {
		// ReadMemory()
		dt, fps = frame.Next()

		loopPhys(dt)
		loopCallbacks(dt)

		loopRenderShadows(dt)
		loopRender(dt)

		loopUI(dt)

		// draw the screen
		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loopCallbacks(dt float32) {
	for i, f := range e.callbacks {
		if f != nil && !f(dt) {
			delete(e.callbacks, i)
		}
	}

	for o := range Objects {
		if o.Callback != nil {
			o.Callback(o, dt)
		}
	}
}

func loopUI(dt float32) {
	e.gfx.Disable(graphicsprovider.DEPTH_TEST)
	e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

	e.ui.Construct(float64(dt))
	e.ui.Draw()

	e.gfx.Disable(graphicsprovider.SCISSOR_TEST)
	e.gfx.Enable(graphicsprovider.DEPTH_TEST)
}

func loopRender(dt float32) {

	w, h := e.window.GetFramebufferSize()

	// clear the screen and reset our viewport
	e.gfx.Viewport(0, 0, int32(w), int32(h))
	e.gfx.ClearColor(0.0, 0.0, 0.0, 10)
	e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

	// make the projection and view matrixes
	perspective := mgl32.Perspective(mgl32.DegToRad(50), float32(w)/float32(h), 1.0, 100.0)
	view := e.camera.GetViewMatrix()

	// render object
	for o, b := range Objects {
		if b && o.Node != nil && !o.Transparent {
			o.Render(perspective, view, e.camera)
		}
	}

	// render scene
	for _, o := range Scene {
		if !o.Transparent {
			o.Render(perspective, view, e.camera)
		}
	}

	// render transparent objects
	for o, b := range Objects {
		if b && o.Node != nil && o.Transparent {
			o.Render(perspective, view, e.camera)
		}
	}

	for _, o := range Scene {
		if o.Transparent {
			o.Render(perspective, view, e.camera)
		}
	}

}

func loopRenderShadows(dt float32) {
	e.render.StartShadowMapping()
	lightCount := e.render.GetActiveLightCount()
	for i := 0; i < lightCount; i++ {
		// get lights with shadow maps
		lightToCast := e.render.ActiveLights[i]
		if lightToCast.ShadowMap == nil {
			continue
		}

		// engine.able the light to cast shadows
		e.render.EnableShadowMappingLight(lightToCast)
		for o, b := range Objects {
			if b && o.Node != nil && o.Shadow {
				e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
			}
		}

		for _, o := range Scene {
			e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
		}
	}
	e.render.EndShadowMapping()
}