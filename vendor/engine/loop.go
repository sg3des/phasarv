package engine

import (
	"engine/frames"
	"fmt"
	"phys"
	"render"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/eweygewey"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

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

			for i, f := range e.callbacks {
				if f != nil && !f(dt) {
					e.callbacks[i] = e.callbacks[len(e.callbacks)-1]
					e.callbacks = e.callbacks[:len(e.callbacks)-1]
					// delete(e.callbacks, i)
				}
			}

			Objects.loopCallbacks(dt)
			phys.NextFrame(dt)
		}
	}()

	for !e.window.ShouldClose() {
		Rdt, Rfps = Rendframe.Next()
		Objects.loopPhysToRender()

		loopRenderFrame(Rdt)
	}
}

func loopRenderFrame(dt float32) {
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
