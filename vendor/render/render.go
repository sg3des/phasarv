package render

import (
	"fmt"
	"log"
	"ui"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/renderer/forward"
)

var (
	render      *forward.ForwardRenderer
	perspective mgl32.Mat4
	view        mgl32.Mat4
	camera      *fizzle.YawPitchCamera
	shadowmap   *fizzle.RenderShader

	Renderables []*Renderable
	Scene       []*Renderable
)

func renderInit(gfx graphicsprovider.GraphicsProvider, w, h int32) {
	render = forward.NewForwardRenderer(gfx)
	render.ChangeResolution(w, h)

	// enable shadow mapping in the renderer
	render.SetupShadowMapRendering()

	var err error
	shadowmap, err = forward.CreateShadowmapGeneratorShader()
	if err != nil {
		log.Fatalln("failed init shadowmap", err)
	}

	// Renderables = make(map[*Renderable]bool)
	// Scene = make(map[*Body]bool)
}

// func SetCamera(c *fizzle.YawPitchCamera) {
// 	camera = c
// }

func DrawFrame(dt float32, fps int) {
	w, h := window.GetFramebufferSize()

	// x, y, _ := camera.GetPosition().Elem()
	// particlesUpdate(float64(dt), mgl32.Vec3{y, x, 0})

	// clear the screen and reset our viewport
	gfx.Viewport(0, 0, int32(w), int32(h))

	// gfx.ClearColor(0.5, 0.5, 0.5, 1)
	gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

	drawParticles(dt)
	drawObjects(mgl32.DegToRad(50), float32(w)/float32(h))

	//render ui
	ui.RenderFrame.Text = []string{fmt.Sprintf("fps: %d dt: %f", fps, dt)}
	renderUI(float64(dt))

	//end frame
	window.SwapBuffers()
	glfw.PollEvents()
}

func drawObjects(fov, aspect float32) {
	// make the projection and view matrixes
	perspective = mgl32.Perspective(fov, aspect, 1.0, 100.0)
	view = camera.GetViewMatrix()

	// log.Println(len(Renderables))

	for i := 0; i < len(Renderables); i++ {
		if Renderables[i].needDestroy {
			DeleteRenderable(i)
			// Renderables[i] = nil
			i--
		}
	}

	// for i, r := range Renderables {
	// 	if i >= len(Renderables) {
	// 		break
	// 	}
	// 	if r.needDestroy {
	// 		DeleteRenderable(i)
	// 	}
	// }

	// render not transparent bodies
	for _, r := range Renderables {
		if r == nil || r.needDestroy {
			// log.Println("DELETEED!!!")
			// DeleteRenderables(i)
			continue
			// r = Renderables[i]
		}
		if !r.Transparent {
			r.render()
		}
	}

	// render not transparent scene objects
	for _, r := range Scene {
		if !r.Transparent {
			r.render()
		}
	}

	// render transparent bodies
	for _, r := range Renderables {
		if r.Transparent {
			r.render()
		}
	}

	// render transparent scene objects
	for _, r := range Scene {
		if r.Transparent {
			r.render()
		}
	}

	// particlesRender(perspective, view)

	renderShadows()
}

func renderShadows() {
	render.StartShadowMapping()
	lightCount := render.GetActiveLightCount()
	for i := 0; i < lightCount; i++ {
		// get lights with shadow maps
		lightToCast := render.ActiveLights[i]
		if lightToCast.ShadowMap == nil {
			continue
		}

		// enable the light to cast shadows
		render.EnableShadowMappingLight(lightToCast)
		for _, r := range Renderables {
			if r.Shadow && r.Body != nil {
				render.DrawRenderableWithShader(r.Body, shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, camera)
			}
		}

		for _, r := range Scene {
			if r.Shadow && r.Body != nil {
				render.DrawRenderableWithShader(r.Body, shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, camera)
			}
		}

	}

	render.EndShadowMapping()
}

func DeleteRenderable(i int) {
	copy(Renderables[i:], Renderables[i+1:])
	Renderables[len(Renderables)-1] = nil // or the zero vRenderableslue of T
	Renderables = Renderables[:len(Renderables)-1]

	// Renderables[i] = nil
	// Renderables[i] = Renderables[len(Renderables)-1]
	// Renderables = Renderables[:len(Renderables)-1]

}

func renderUI(dt float64) {
	gfx.Disable(graphicsprovider.DEPTH_TEST)
	gfx.Enable(graphicsprovider.SCISSOR_TEST)

	ui.Draw(dt)

	gfx.Disable(graphicsprovider.SCISSOR_TEST)
	gfx.Enable(graphicsprovider.DEPTH_TEST)
}
