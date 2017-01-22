package engine

import (
	"engine/frames"
	"fmt"
	"log"
	"phys"
	"render"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

func LoopPlay() {
	frame := frames.NewFrame()
	var dt float32
	var fps int
	for {
		dt, fps = frame.Next()
		log.Println(dt, fps)

		if UI.PhysFrame != nil {
			UI.PhysFrame.Text = []string{fmt.Sprintf("fps: %d dt: %f", fps, dt)}
		}

		for i, f := range e.callbacks {
			if f != nil && !f(dt) {
				e.callbacks[i] = e.callbacks[len(e.callbacks)-1]
				e.callbacks = e.callbacks[:len(e.callbacks)-1]
				// delete(e.callbacks, i)
			}
		}

		phys.NextFrame(dt)
		Objects.loopCallbacks(dt)
	}
}

func LoopRender() {
	frame := frames.NewFrame()
	var dt float32
	var fps int
	for !e.window.ShouldClose() {

		Objects.loopPhysToRender()

		dt, fps = frame.Next()
		w, h := e.window.GetFramebufferSize()

		// clear the screen and reset our viewport
		e.gfx.Viewport(0, 0, int32(w), int32(h))

		// e.gfx.ClearColor(0.5, 0.5, 0.5, 1)
		e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		render.NextFrame(mgl32.DegToRad(50), float32(w)/float32(h))

		// log.Println(dt)
		UI.RenderFrame.Text = []string{fmt.Sprintf("fps: %d dt: %f", fps, dt)}
		uiConstruct(dt)

		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}
