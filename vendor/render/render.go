package render

import (
	"log"

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

	Renderables map[*Renderable]bool
	Scene       []*Renderable
)

func Init(gfx graphicsprovider.GraphicsProvider, w, h int32) {
	render = forward.NewForwardRenderer(gfx)
	render.ChangeResolution(w, h)

	// enable shadow mapping in the renderer
	render.SetupShadowMapRendering()

	var err error
	shadowmap, err = forward.CreateShadowmapGeneratorShader()
	if err != nil {
		log.Fatalln("failed init shadowmap", err)
	}

	Renderables = make(map[*Renderable]bool)
	// Scene = make(map[*Body]bool)
}

func SetCamera(c *fizzle.YawPitchCamera) {
	camera = c
}

func NextFrame(fov, aspect float32) {
	// make the projection and view matrixes
	perspective = mgl32.Perspective(fov, aspect, 1.0, 100.0)
	view = camera.GetViewMatrix()

	// renderShadows()

	// render not transparent bodies
	for r := range Renderables {
		if !r.Transparent {
			r.Render()
		}
	}

	// render not transparent scene objects
	for _, r := range Scene {
		if !r.Transparent {
			r.Render()
		}
	}

	// render transparent bodies
	for r := range Renderables {
		if r.Transparent {
			r.Render()
		}
	}

	// render transparent scene objects
	for _, r := range Scene {
		if r.Transparent {
			r.Render()
		}
	}
}

func renderShadows() {
	render.StartShadowMapping()
	lightCount := render.GetActiveLightCount()
	log.Println(lightCount)
	for i := 0; i < lightCount; i++ {
		// get lights with shadow maps
		lightToCast := render.ActiveLights[i]
		if lightToCast.ShadowMap == nil {
			continue
		}

		// enable the light to cast shadows
		render.EnableShadowMappingLight(lightToCast)
		for r := range Renderables {
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
