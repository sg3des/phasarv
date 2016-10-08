package engine

import (
	"assets"
	"engine/frames"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"
	"github.com/tbogdala/fizzle/input/glfwinput"
	"github.com/tbogdala/fizzle/renderer/forward"
)

var engine struct {
	Window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider
	render *forward.ForwardRenderer

	shadowmap *fizzle.RenderShader

	Camera *fizzle.YawPitchCamera

	keyboard *glfwinput.KeyboardModel

	callbacks []func(float32)
}

//init main function create window and initialize opengl,render engine.c.
func init() {
	runtime.LockOSThread()

	var err error

	engine.Window, err = initWindow(1024, 768, "phasarv")
	if err != nil {
		log.Fatalf("engine: failed initialize window, reason: %s", err)
	}

	engine.gfx, err = initOpenGL()
	if err != nil {
		log.Fatalf("engine: failed initialize opengl, reason: %s", err)
	}

	engine.keyboard = initKeyboard(engine.Window)
	engine.render = initRender(engine.gfx)

	err = assets.LoadAssets("assets/textures", "assets/shaders", "assets/models")
	if err != nil {
		log.Fatalln("failed load assets", err)
	}

	engine.shadowmap, err = forward.CreateShadowmapGeneratorShader()
	if err != nil {
		log.Fatalf("failed generate shadow map, reason: %s", err)
	}

	// engine.Camera = fizzle.NewYawPitchCamera(mgl32.Vec3{0, 0, 10})
}

//initWindow create window and set some opengl flags
func initWindow(w, h int, title string) (*glfw.Window, error) {
	// GLFW must be initialized before it's called
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 1)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	window, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	// v-sync
	glfw.SwapInterval(0)
	return window, nil
}

//initOpenGL initialize openGL
func initOpenGL() (graphicsprovider.GraphicsProvider, error) {
	gfx, err := opengl.InitOpenGL()
	if err != nil {
		return nil, err
	}
	fizzle.SetGraphics(gfx)
	return gfx, err
}

func initKeyboard(w *glfw.Window) *glfwinput.KeyboardModel {
	kbModel := glfwinput.NewKeyboardModel(w)
	kbModel.BindTrigger(glfw.KeyEscape, func() {
		w.SetShouldClose(true)
	})
	kbModel.SetupCallbacks()

	return kbModel
}

func initRender(gfx graphicsprovider.GraphicsProvider) *forward.ForwardRenderer {
	// create a new renderer
	render := forward.NewForwardRenderer(gfx)
	render.ChangeResolution(1024, 768)

	// engine.able shadow mapping in the renderer
	render.SetupShadowMapRendering()

	// set some OpenGL flags
	gfx.Enable(graphicsprovider.CULL_FACE)
	gfx.Enable(graphicsprovider.DEPTH_TEST)
	gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
	gfx.Enable(graphicsprovider.TEXTURE_2D)
	gfx.Enable(graphicsprovider.BLEND)
	gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)

	return render
}

//AddCallback add  engine.ch frame callback with deltaT as one argument
func AddCallback(funcs ...func(float32)) {
	for _, f := range funcs {
		engine.callbacks = append(engine.callbacks, f)
	}
}

func SetMouseCallback(f func(*glfw.Window, glfw.MouseButton, glfw.Action, glfw.ModifierKey)) {
	engine.Window.SetMouseButtonCallback(f)
}

func CursorPosition() (float32, float32) {
	x, y := engine.Window.GetCursorPos()
	return float32(x), float32(y)
}

func WindowSize() (float32, float32) {
	w, h := engine.Window.GetSize()
	return float32(w), float32(h)
}

func renderShadows() {
	engine.render.StartShadowMapping()
	lightCount := engine.render.GetActiveLightCount()
	for i := 0; i < lightCount; i++ {
		// get lights with shadow maps
		lightToCast := engine.render.ActiveLights[i]
		if lightToCast.ShadowMap == nil {
			continue
		}

		// engine.able the light to cast shadows
		engine.render.EnableShadowMappingLight(lightToCast)
		for o, b := range Objects {
			if b && o.Node != nil && o.Shadow {
				engine.render.DrawRenderableWithShader(o.Node, engine.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, engine.Camera)
			}
		}
	}
	engine.render.EndShadowMapping()
}

func Loop() {
	frame := frames.NewFrame()
	InitPhys(0.3)

	// loop until something told the Window that it should close
	for !engine.Window.ShouldClose() {
		frame.Next()
		frame.FPS()
		dt := frame.DT()
		physRender(dt)

		for _, f := range engine.callbacks {
			f(dt)
		}

		for o := range Objects {
			if o.Callback != nil {
				o.Callback(o, dt)
			}
		}

		renderShadows()

		w, h := engine.Window.GetFramebufferSize()

		// clear the screen and reset our viewport
		engine.gfx.Viewport(0, 0, int32(w), int32(h))
		// gfx.ClearColor(0.1, 0.1, 0.1, 0.1)
		engine.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		// make the projection and view matrixes
		perspective := mgl32.Perspective(mgl32.DegToRad(50), float32(w)/float32(h), 1.0, 100.0)
		view := engine.Camera.GetViewMatrix()

		// render object
		for o, b := range Objects {
			if b && o.Node != nil && !o.Transparent {
				o.Render(perspective, view, engine.Camera)
			}
		}

		// render transparent objects
		for o, b := range Objects {
			if b && o.Node != nil && o.Transparent {
				o.Render(perspective, view, engine.Camera)
			}
		}

		// draw the screen
		engine.Window.SwapBuffers()
		glfw.PollEvents()
	}
}
