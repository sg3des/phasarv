package engine

import (
	"assets"
	"fmt"
	"log"
	"phys"
	"render"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"

	"github.com/tbogdala/eweygewey"
	"github.com/tbogdala/eweygewey/glfwinput"
)

var (
	err error
	e   = &engine{}
	mem runtime.MemStats
)

type engine struct {
	window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider
	// render *forward.ForwardRenderer

	ui *eweygewey.Manager

	// shadowmap *fizzle.RenderShader

	// camera *fizzle.YawPitchCamera

	callbacks map[int]func(float32) bool
}

//Init main function create window and initialize opengl,render engine
func Init(userfunc func()) {
	// runtime.LockOSThread()

	if err := e.initWindow(1024, 768, "phasarv"); err != nil {
		log.Fatalf("engine: failed initialize window, reason: %s", err)
	}

	if err := e.initOpenGL(); err != nil {
		log.Fatalf("engine: failed initialize opengl, reason: %s", err)
	}

	if err := assets.LoadAssets("assets/textures", "assets/shaders", "assets/models"); err != nil {
		log.Fatalf("failed load assets, reason: %s", err)
	}

	// e.initRender()

	// if err := e.initShadowmap(); err != nil {
	// 	log.Fatalf("failed generate shadow map, reason: %s", err)
	// }

	render.Init(e.gfx, 1024, 768)
	render.SetCamera(fizzle.NewYawPitchCamera(mgl32.Vec3{0, 0, 10}))

	phys.Init()

	if err := e.initUI(); err != nil {
		log.Fatalf("failed initialise user interface, %s", err)
	}

	e.callbacks = make(map[int]func(float32) bool)

	// e.camera = fizzle.NewYawPitchCamera(mgl32.Vec3{0, 0, 10})

	userfunc()

	Loop()
}

//initWindow create window and set some opengl flags
func (e *engine) initWindow(w, h int, title string) error {
	// GLFW must be initialized before it's called
	if err := glfw.Init(); err != nil {
		return err
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 4)
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

	e.gfx.Enable(graphicsprovider.TEXTURE_2D)
	e.gfx.Enable(graphicsprovider.BLEND)
	// e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

	// e.gfx.Enable(graphicsprovider.CULL_FACE)
	// e.gfx.Enable(graphicsprovider.DEPTH_TEST)
	// e.gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
	// e.gfx.Enable(graphicsprovider.TEXTURE_2D)
	// e.gfx.Enable(graphicsprovider.BLEND)

	e.gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)
	e.gfx.Enable(graphicsprovider.SAMPLE_ALPHA_TO_COVERAGE)

	// e.gfx.Enable(graphicsprovider.EQUAL)
	// e.gfx.BlendFunc(graphicsprovider.ONE_MINUS_DST_ALPHA, graphicsprovider.DST_ALPHA)
	// e.gfx.BlendFunc(graphicsprovider.ONE, graphicsprovider.ONE)

	// glBlendFunc(GL_ONE_MINUS_DST_ALPHA, GL_DST_ALPHA)

	return nil
}

// func (e *engine) initRender() {
// 	// e.render = forward.NewForwardRenderer(e.gfx)
// 	// e.render.ChangeResolution(1024, 768)

// 	// // engine.able shadow mapping in the renderer
// 	// e.render.SetupShadowMapRendering()

// 	// set some OpenGL flags
// 	// e.gfx.Enable(graphicsprovider.CULL_FACE)
// 	// e.gfx.Enable(graphicsprovider.DEPTH_TEST)
// 	e.gfx.Enable(graphicsprovider.TEXTURE_2D)
// 	e.gfx.Enable(graphicsprovider.BLEND)
// 	// e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

// 	// e.gfx.Enable(graphicsprovider.CULL_FACE)
// 	// e.gfx.Enable(graphicsprovider.DEPTH_TEST)
// 	// e.gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
// 	// e.gfx.Enable(graphicsprovider.TEXTURE_2D)
// 	// e.gfx.Enable(graphicsprovider.BLEND)

// 	e.gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)
// 	e.gfx.Enable(graphicsprovider.SAMPLE_ALPHA_TO_COVERAGE)

// 	// e.gfx.Enable(graphicsprovider.EQUAL)
// 	// e.gfx.BlendFunc(graphicsprovider.ONE_MINUS_DST_ALPHA, graphicsprovider.DST_ALPHA)
// 	// e.gfx.BlendFunc(graphicsprovider.ONE, graphicsprovider.ONE)

// 	// glBlendFunc(GL_ONE_MINUS_DST_ALPHA, GL_DST_ALPHA)
// }

// func (e *engine) initShadowmap() error {
// 	e.shadowmap, err = forward.CreateShadowmapGeneratorShader()
// 	return err
// }

func (e *engine) initUI() error {
	fontScale := 14
	fontFilepath := "assets/fonts/Roboto-Bold.ttf"
	fontGlyphs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890., :[]{}\\|<>;\"'~`?/-+_=()*&^%$#@!"

	// create and initialize the gui Manager
	e.ui = eweygewey.NewManager(e.gfx)
	if err := e.ui.Initialize(eweygewey.VertShader330, eweygewey.FragShader330, 1024, 768, 768); err != nil {
		return fmt.Errorf("Failed to initialize the user interface! reason %s", err)
	}

	glfwinput.SetInputHandlers(e.ui, e.window)

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
		e.callbacks[len(e.callbacks)] = f
	}
}

//SetMouseCallback set function  callback each frame
func SetMouseCallback(f func(*glfw.Window, glfw.MouseButton, glfw.Action, glfw.ModifierKey)) {
	e.window.SetMouseButtonCallback(f)
}

//SetKeyCallback set function  callback each frame
func SetKeyCallback(f func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)) {
	e.window.SetKeyCallback(f)
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
