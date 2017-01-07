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
		loopPhysToRender()
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

//loopPhysToRender update renderable position and rotation for dynamical objects
func loopPhysToRender() {
	for o := range Objects {
		if o.needDestroy {
			o.renderable.Destroy()
			delete(Objects, o)
			continue
		}
		if o.renderable.Body == nil {
			continue
		}
		// update position
		o.renderable.Body.Location = o.PositionVec3()

		// update rotation
		ang := o.Rotation()
		//if rollAngle exist then need roll renderable object
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
		if o.needDestroy {
			continue
		}
		for _, f := range o.Callbacks {

			f(dt)
		}
	}
}
