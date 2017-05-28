package engine

import (
	"engine/frames"
	"phys"
	"render"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
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
	// var mem runtime.MemStats

	frame := frames.NewFrame()
	for !window.ShouldClose() {
		// runtime.ReadMemStats(&mem)

		// log.Println(mem.Alloc/1024, mem.TotalAlloc/1024, mem.HeapAlloc/1024, mem.HeapSys/1024)

		dt, fps := frame.Next()

		frameLogic(dt)
		frameRender(dt, fps)
		glfw.PollEvents()
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
