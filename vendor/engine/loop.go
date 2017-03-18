package engine

import (
	"engine/frames"
	"phys"
	"render"
	"time"
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
	for !window.ShouldClose() {
		dt, fps := frame.Next()

		frameLogic(dt)
		frameRender(dt, fps)
	}
}

func frameLogic(dt float32) {
	for i, f := range callbacks {
		if f != nil && !f(dt) {
			callbacks[i] = callbacks[len(callbacks)-1]
			callbacks = callbacks[:len(callbacks)-1]
			// delete(e.callbacks, i)
		}
	}

	phys.NextFrame(dt)
	Objects.loopCallbacks(dt)
}

func frameRender(dt float32, fps int) {
	Objects.loopPhysToRender()
	render.DrawFrame(dt, fps)
}
