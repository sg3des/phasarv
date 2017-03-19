package render

import (
	"fmt"
	"runtime"
	"ui"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"
)

var (
	window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider
)

//NewWindow initialize all sub steps for creation window with renderable content
func NewWindow(width, height int, title string) (*glfw.Window, error) {
	runtime.LockOSThread()

	if err := initWindow(width, height, title); err != nil {
		return nil, fmt.Errorf("engine: failed initialize window, reason: %s", err)
	}

	if err := initOpenGL(); err != nil {
		return nil, fmt.Errorf("engine: failed initialize opengl, reason: %s", err)
	}

	renderInit(gfx, int32(width), int32(height))
	NewCamera(mgl32.Vec3{0, 0, 10})

	if err := ui.Init(gfx, window); err != nil {
		return nil, fmt.Errorf("failed initialize UI, reason: %s", err)
	}

	return window, nil
}

//initWindow create window and set some opengl flags
func initWindow(w, h int, title string) error {
	// GLFW must be initialized before it's called
	err := glfw.Init()
	if err != nil {
		return err
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	window, err = glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	// v-sync
	glfw.SwapInterval(1)
	return nil
}

//initOpenGL initialize openGL
func initOpenGL() error {
	var err error
	gfx, err = opengl.InitOpenGL()
	if err != nil {
		return err
	}

	fizzle.SetGraphics(gfx)

	gfx.Enable(graphicsprovider.TEXTURE_2D)
	gfx.Enable(graphicsprovider.BLEND)
	// e.gfx.Enable(graphicsprovider.SCISSOR_TEST)

	// e.gfx.Enable(graphicsprovider.CULL_FACE)
	// e.gfx.Enable(graphicsprovider.DEPTH_TEST)
	// e.gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
	// e.gfx.Enable(graphicsprovider.TEXTURE_2D)
	// e.gfx.Enable(graphicsprovider.BLEND)

	gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)
	gfx.Enable(graphicsprovider.SAMPLE_ALPHA_TO_COVERAGE)

	// e.gfx.Enable(graphicsprovider.EQUAL)
	// e.gfx.BlendFunc(graphicsprovider.ONE_MINUS_DST_ALPHA, graphicsprovider.DST_ALPHA)
	// e.gfx.BlendFunc(graphicsprovider.ONE, graphicsprovider.ONE)

	// glBlendFunc(GL_ONE_MINUS_DST_ALPHA, GL_DST_ALPHA)

	return nil
}