package engine

import (
	"engine/frames"
	"fmt"
	"phys"
	"render"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

func LoopServer() {
	frame := frames.NewFrame()
	ticker := time.NewTicker(time.Second / 60)
	for _ = range ticker.C {
		dt, _ := frame.Next()

		frameLogic(dt)
	}
}

func LoopRender() {
	frame := frames.NewFrame()
	for !e.window.ShouldClose() {
		dt, fps := frame.Next()

		frameLogic(dt)
		frameRender(dt, fps)
	}
}

func frameLogic(dt float32) {
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

func frameRender(dt float32, fps int) {
	Objects.loopPhysToRender()
	w, h := e.window.GetFramebufferSize()

	// clear the screen and reset our viewport
	e.gfx.Viewport(0, 0, int32(w), int32(h))

	// e.gfx.ClearColor(0.5, 0.5, 0.5, 1)
	e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

	render.NextFrame(mgl32.DegToRad(50), float32(w)/float32(h))

	//render ui
	UI.RenderFrame.Text = []string{fmt.Sprintf("fps: %d dt: %f", fps, dt)}
	uiConstruct(float64(dt))

	//end frame
	e.window.SwapBuffers()
	glfw.PollEvents()
}
