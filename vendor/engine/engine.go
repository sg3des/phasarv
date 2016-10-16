package engine

import (
	"assets"
	"engine/frames"
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"
	"github.com/tbogdala/fizzle/input/glfwinput"
	"github.com/tbogdala/fizzle/renderer/forward"

	gui "github.com/tbogdala/eweygewey"
	guiinput "github.com/tbogdala/eweygewey/glfwinput"
)

var (
	err error
	e   = &engine{}
)

type engine struct {
	window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider
	render *forward.ForwardRenderer

	ui *gui.Manager

	shadowmap *fizzle.RenderShader

	camera *fizzle.YawPitchCamera

	keyboard *glfwinput.KeyboardModel

	callbacks []func(float32) bool
}

//init main function create window and initialize opengl,render engine.c.
func init() {
	runtime.LockOSThread()

	if err := e.initWindow(1024, 768, "phasarv"); err != nil {
		log.Fatalf("engine: failed initialize window, reason: %s", err)
	}

	if err := e.initOpenGL(); err != nil {
		log.Fatalf("engine: failed initialize opengl, reason: %s", err)
	}

	e.initKeyboard()
	e.initRender()

	if err := assets.LoadAssets("assets/textures", "assets/shaders", "assets/models"); err != nil {
		log.Fatalf("failed load assets, reason: %s", err)
	}

	if err := e.initShadowmap(); err != nil {
		log.Fatalf("failed generate shadow map, reason: %s", err)
	}

	if err := e.initUI(); err != nil {
		log.Fatalf("failed initialise user interface, %s", err)
	}

	// engine.Camera = fizzle.NewYawPitchCamera(mgl32.Vec3{0, 0, 10})
}

//initWindow create window and set some opengl flags
func (e *engine) initWindow(w, h int, title string) error {
	// GLFW must be initialized before it's called
	if err := glfw.Init(); err != nil {
		return err
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 1)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	e.window, err = glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return err
	}
	e.window.MakeContextCurrent()

	// v-sync
	glfw.SwapInterval(0)
	return nil
}

//initOpenGL initialize openGL
func (e *engine) initOpenGL() error {
	e.gfx, err = opengl.InitOpenGL()
	if err != nil {
		return err
	}

	fizzle.SetGraphics(e.gfx)
	return nil
}

func (e *engine) initKeyboard() {
	e.keyboard = glfwinput.NewKeyboardModel(e.window)
	e.keyboard.BindTrigger(glfw.KeyEscape, func() {
		e.window.SetShouldClose(true)
	})

	//enable/disable limit 60fps
	var limitFPS bool
	e.keyboard.BindTrigger(glfw.KeyS, func() {
		if limitFPS {
			glfw.SwapInterval(0)
			limitFPS = false
		} else {
			glfw.SwapInterval(1)
			limitFPS = true
		}

	})
	e.keyboard.SetupCallbacks()

}

func (e *engine) initRender() {
	// create a new renderer
	e.render = forward.NewForwardRenderer(e.gfx)
	e.render.ChangeResolution(1024, 768)

	// engine.able shadow mapping in the renderer
	e.render.SetupShadowMapRendering()

	// set some OpenGL flags
	e.gfx.Enable(graphicsprovider.CULL_FACE)
	e.gfx.Enable(graphicsprovider.DEPTH_TEST)
	e.gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
	e.gfx.Enable(graphicsprovider.TEXTURE_2D)
	e.gfx.Enable(graphicsprovider.BLEND)
	e.gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.
		ONE_MINUS_SRC_ALPHA)
	e.gfx.Enable(graphicsprovider.SAMPLE_ALPHA_TO_COVERAGE)

	e.gfx.Enable(graphicsprovider.EQUAL)
	// e.gfx.BlendFunc(graphicsprovider.ONE_MINUS_DST_ALPHA, graphicsprovider.DST_ALPHA)
	// gfx.BlendFunc(graphicsprovider.ONE, graphicsprovider.ONE)

	// glBlendFunc(GL_ONE_MINUS_DST_ALPHA, GL_DST_ALPHA)
}

func (e *engine) initShadowmap() error {
	e.shadowmap, err = forward.CreateShadowmapGeneratorShader()
	return err
}

func (e *engine) initUI() error {
	fontScale := 14
	fontFilepath := "assets/fonts/Roboto-Bold.ttf"
	fontGlyphs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890., :[]{}\\|<>;\"'~`?/-+_=()*&^%$#@!"

	// create and initialize the gui Manager
	e.ui = gui.NewManager(e.gfx)
	if err := e.ui.Initialize(gui.VertShader330, gui.FragShader330, 1024, 768, 768); err != nil {
		return fmt.Errorf("Failed to initialize the user interface! reason %s", err)
	}
	guiinput.SetInputHandlers(e.ui, e.window)

	// load a font
	_, err = e.ui.NewFont("Default", fontFilepath, fontScale, fontGlyphs)
	if err != nil {
		return fmt.Errorf("Failed to load the font file! reason: %s", err)
	}

	return nil
}

//AddCallback add  engine.ch frame callback with deltaT as one argument
func AddCallback(funcs ...func(float32) bool) {
	for _, f := range funcs {
		e.callbacks = append(e.callbacks, f)
	}
}

//SetMouseCallback set function  callback each frame
func SetMouseCallback(f func(*glfw.Window, glfw.MouseButton, glfw.Action, glfw.ModifierKey)) {
	e.window.SetMouseButtonCallback(f)
}

//CursorPosition returned cursor position
func CursorPosition() (float32, float32) {
	x, y := e.window.GetCursorPos()
	return float32(x), float32(y)
}

//WindowSize returned current size of window [type: float32]
func WindowSize() (float32, float32) {
	w, h := e.window.GetSize()
	return float32(w), float32(h)
}

//Loop is function where:  poll events, calculate physics, rendered objects and effects. This is infinite loop, and should be run late from user space
func Loop() {
	frame := frames.NewFrame()
	InitPhys(0.3)

	var dt float32
	var fps int

	wnd := e.ui.NewWindow("test", 0.79, 0.99, 0.2, 0, func(wng *gui.Window) {
		wng.Text(fmt.Sprintf("fps: %d dt: %f", fps, dt))
	})
	wnd.ShowTitleBar = false
	wnd.IsMoveable = false
	wnd.AutoAdjustHeight = true

	// loop until something told the Window that it should close
	for !e.window.ShouldClose() {
		dt, fps = frame.Next()

		// log.Println(len(e.callbacks))

		loopRenderPhys(dt)

		for i, f := range e.callbacks {
			if f != nil && !f(dt) {
				copy(e.callbacks[i:], e.callbacks[i+1:])       // shift
				e.callbacks[len(e.callbacks)-1] = nil          // remove reference
				e.callbacks = e.callbacks[:len(e.callbacks)-1] // reslice
			}
			// e.callbacks[i] = nil

		}

		/*	for i, f := range e.callbacks {
				if f == nil {
					copy(e.callbacks[i:], e.callbacks[i+1:])       // shift
					e.callbacks[len(e.callbacks)-1] = nil          // remove reference
					e.callbacks = e.callbacks[:len(e.callbacks)-1] // reslice
				}
			}
		*/
		for o := range Objects {
			if o.Callback != nil {
				o.Callback(o, dt)
			}
		}

		loopRenderParticles(dt)
		loopRenderShadows()

		w, h := e.window.GetFramebufferSize()

		// clear the screen and reset our viewport
		e.gfx.Viewport(0, 0, int32(w), int32(h))
		// gfx.ClearColor(0.1, 0.1, 0.1, 0.1)
		e.gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		// make the projection and view matrixes
		perspective := mgl32.Perspective(mgl32.DegToRad(50), float32(w)/float32(h), 1.0, 100.0)
		view := e.camera.GetViewMatrix()

		// render object
		for o, b := range Objects {
			if b && o.Node != nil && !o.Transparent {
				o.Render(perspective, view, e.camera)
			}
		}

		// render scene
		for _, o := range Scene {
			if !o.Transparent {
				o.Render(perspective, view, e.camera)
			}
		}

		// render transparent objects
		for o, b := range Objects {
			if b && o.Node != nil && o.Transparent {
				o.Render(perspective, view, e.camera)
			}
		}

		for _, o := range Scene {
			if o.Transparent {
				o.Render(perspective, view, e.camera)
			}
		}

		// draw the user interface
		e.gfx.Disable(graphicsprovider.DEPTH_TEST)
		e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

		e.ui.Construct(float64(dt))
		e.ui.Draw()

		e.gfx.Disable(graphicsprovider.SCISSOR_TEST)
		e.gfx.Enable(graphicsprovider.DEPTH_TEST)

		// draw the screen
		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loopRenderShadows() {
	e.render.StartShadowMapping()
	lightCount := e.render.GetActiveLightCount()
	for i := 0; i < lightCount; i++ {
		// get lights with shadow maps
		lightToCast := e.render.ActiveLights[i]
		if lightToCast.ShadowMap == nil {
			continue
		}

		// engine.able the light to cast shadows
		e.render.EnableShadowMappingLight(lightToCast)
		for o, b := range Objects {
			if b && o.Node != nil && o.Shadow {
				e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
			}
		}

		for _, o := range Scene {
			e.render.DrawRenderableWithShader(o.Node, e.shadowmap, nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, e.camera)
		}
	}
	e.render.EndShadowMapping()
}
