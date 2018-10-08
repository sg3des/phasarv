package render

import (
	"assets"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sg3des/fizzgui"
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

func renderInit(w, h int32) {
	fizzle.SetGraphics(gfx)

	render = forward.NewForwardRenderer(gfx)
	render.ChangeResolution(w, h)

	// enable shadow mapping in the renderer
	render.SetupShadowMapRendering()

	var err error
	shadowmap, err = forward.CreateShadowmapGeneratorShader()
	if err != nil {
		log.Fatalln("failed init shadowmap", err)
	}
}

func DrawFrame(dt float32) {
	w, h := window.GetFramebufferSize()

	// clear the screen and reset our viewport
	gfx.Viewport(0, 0, int32(w), int32(h))
	// gfx.ClearColor(0.4, 0.4, 0.4, 1)
	// gfx.ClearColor(0, 0, 0, 0.1)
	gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

	drawParticles(dt)
	drawObjects(mgl32.DegToRad(50), float32(w)/float32(h))

	renderUI()

	//end frame
	window.SwapBuffers()
}

func drawObjects(fov, aspect float32) {
	// make the projection and view matrixes
	perspective = mgl32.Perspective(fov, aspect, 1.0, 100.0)
	view = camera.GetViewMatrix()

	Renderables = DeleteRenderables(Renderables)
	Scene = DeleteRenderables(Scene)

	// render not transparent bodies
	for _, r := range Renderables {
		// if r == nil || r.needDestroy {
		// 	continue
		// }
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

func DeleteRenderables(a []*Renderable) []*Renderable {
	for i := 0; i < len(a); i++ {
		if a[i].needDestroy {

			a[i] = a[len(a)-1]
			a[len(a)-1] = nil
			a = a[:len(a)-1]

			i--
		}
	}
	return a
}

// func DeleteRenderable(a []*Renderable, i int) []*Renderable {
// 	copy(a[i:], a[i+1:])
// 	a[len(a)-1] = nil
// 	return a[:len(a)-1]
// }

//
// UI
//

func InitUI(uiDir string) error {
	fizzgui.DefaultPadding = fizzgui.Offset{10, 10, 10, 10}
	fizzgui.BorderColor = fizzgui.Color(70, 130, 220, 50)
	fizzgui.BorderColorHiglight = fizzgui.Color(140, 200, 200, 250)

	fizzgui.BGColorBtn = fizzgui.Color(15, 45, 55, 255)
	fizzgui.BGColorBtnHover = fizzgui.Color(25, 55, 65, 255)
	fizzgui.BGColorHighlight = fizzgui.Color(35, 65, 75, 255)

	err := fizzgui.Init(window, gfx)
	if err != nil {
		return err
	}

	fizzgui.DefaultContainerStyle = fizzgui.NewStyle(mgl32.Vec4{}, fizzgui.Color(60, 70, 90, 225), fizzgui.BorderColor, 2)

	return loadFonts()
}

func loadFonts() (err error) {
	for _, font := range assets.GetFonts() {
		font.Font, err = fizzgui.NewFont(font.Name, font.Path, font.Size, fizzgui.FontGlyphs)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderUI() {
	// gfx.Disable(graphicsprovider.DEPTH_TEST)
	// gfx.Enable(graphicsprovider.SCISSOR_TEST)

	fizzgui.Construct()

	// gfx.Disable(graphicsprovider.SCISSOR_TEST)
	// gfx.Enable(graphicsprovider.DEPTH_TEST)
}
