package engine

import (
	"assets"
	"engine/frames"
	"log"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/graphicsprovider/opengl"
	"github.com/tbogdala/fizzle/input/glfwinput"
	"github.com/tbogdala/fizzle/renderer/forward"
)

const (
	width  = 1024
	height = 768
	fov    = 50.0
)

var (
	Window *glfw.Window
	gfx    graphicsprovider.GraphicsProvider

	Camera *fizzle.YawPitchCamera

	render *forward.ForwardRenderer

	callbacks []func(float32)
)

func AddCallback(funcs ...func(float32)) {
	for _, f := range funcs {
		callbacks = append(callbacks, f)
	}
}

func NewWindow() {
	initGraphics("go", width, height)

	// set the callback function for key glfwinput
	kbModel := glfwinput.NewKeyboardModel(Window)
	kbModel.BindTrigger(glfw.KeyEscape, setShouldClose)
	kbModel.SetupCallbacks()

	// create a new renderer
	render = forward.NewForwardRenderer(gfx)
	// defer r.Destroy()

	if err := assets.LoadAssets("assets/textures", "assets/shaders", "assets/models"); err != nil {
		panic(err)
	}

	// enable shadow mapping in the renderer
	render.SetupShadowMapRendering()

	// set some OpenGL flags
	gfx.Enable(graphicsprovider.CULL_FACE)
	gfx.Enable(graphicsprovider.DEPTH_TEST)
	gfx.Enable(graphicsprovider.PROGRAM_POINT_SIZE)
	gfx.Enable(graphicsprovider.TEXTURE_2D)
	gfx.Enable(graphicsprovider.BLEND)
	gfx.BlendFunc(graphicsprovider.SRC_ALPHA, graphicsprovider.ONE_MINUS_SRC_ALPHA)

	//hide cursor
	// Window.SetCursor(glfw.CreateCursor(assets.GetImage("assets/textures/hide.png"), 0, 0))
}

// initGraphics creates an OpenGL window and initializes the required graphics libraries.
// It will either succeed or panic.
func initGraphics(title string, w int, h int) (*glfw.Window, graphicsprovider.GraphicsProvider) {
	// GLFW must be initialized before it's called
	if err := glfw.Init(); err != nil {
		log.Fatal("Can't init glfw!", err)
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 1)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	var err error
	Window, err = glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		log.Fatal("Failed to create the main window! " + err.Error())
	}
	Window.MakeContextCurrent()

	// v-sync
	glfw.SwapInterval(0)

	// initialize OpenGL
	gfx, err = opengl.InitOpenGL()
	if err != nil {
		panic("Failed to initialize OpenGL! " + err.Error())
	}
	fizzle.SetGraphics(gfx)

	return Window, gfx
}

// setShouldClose should be called to close the window and kill the app.
func setShouldClose() {
	Window.SetShouldClose(true)
}

//MainLoop - is main loop ^^
func MainLoop() {
	// loop until something told the Window that it should close
	for frame := frames.NewFrame(); !Window.ShouldClose(); frame.Next() {
		frame.FPS()

		dt := frame.DT()

		physRender(dt)

		for _, f := range callbacks {
			f(dt)
		}

		for o := range Objects {
			if o.Callback != nil {
				o.Callback(o, dt)
			}
		}

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
			for o, b := range Objects {
				if b && o.Node != nil && o.Shadow {
					render.DrawRenderableWithShader(o.Node, assets.Shaders["shadowmap_generator"], nil, lightToCast.ShadowMap.Projection, lightToCast.ShadowMap.View, Camera)
				}
			}
		}
		render.EndShadowMapping()

		// clear the screen and reset our viewport
		gfx.Viewport(0, 0, int32(width), int32(height))
		// gfx.ClearColor(0.1, 0.1, 0.1, 0.1)
		gfx.Clear(graphicsprovider.COLOR_BUFFER_BIT | graphicsprovider.DEPTH_BUFFER_BIT)

		// make the projection and view matrixes
		perspective := mgl32.Perspective(mgl32.DegToRad(fov), float32(width)/float32(height), 1.0, 100.0)
		view := Camera.GetViewMatrix()

		// draw the stuff
		for o, b := range Objects {
			if b && o.Node != nil && o.Transparent == false {
				// log.Println(o.Name, o.Node)
				render.DrawRenderable(o.Node, nil, perspective, view, Camera)
			}
		}

		// render childs
		for o, b := range Objects {
			if b && o.Node != nil {

				for _, child := range o.ArtStatic {
					child.Art.Location = o.PositionVec3().Add(child.LocalPosition)
					if child.Line {
						render.DrawLines(child.Art, child.Art.Core.Shader, nil, perspective, view, Camera)
					} else {
						render.DrawRenderable(child.Art, nil, perspective, view, Camera)
					}
				}

				for _, child := range o.ArtRotate {
					// log.Println(child.Name, child.Line)
					// child.Art.Location = o.PositionVec3()
					// child.Art.Location = o.PositionVec3().Add(mgl32.Vec3{0, 0, child.LocalPosition.Z()})

					xF, yF := o.VectorForward(child.LocalPosition.X())
					xS, yS := o.VectorSide(child.LocalPosition.Y())
					child.Art.Location = o.PositionVec3().Add(mgl32.Vec3{xF, yF, child.LocalPosition.Z()}).Add(mgl32.Vec3{xS, yS})
					child.Art.LocalRotation = mgl32.AnglesToQuat(0, 0, o.Rotation(), 1)

					if child.Line {
						render.DrawLines(child.Art, child.Art.Core.Shader, nil, perspective, view, Camera)
					} else {
						render.DrawRenderable(child.Art, nil, perspective, view, Camera)
					}
				}

			}
		}

		// render transparent objects
		for o, b := range Objects {
			if b && o.Node != nil && o.Transparent == true {
				render.DrawRenderable(o.Node, nil, perspective, view, Camera)
			}
		}

		for o, b := range Lines {
			if b && o != nil {
				render.DrawLines(o, o.Core.Shader, nil, perspective, view, Camera)
			}
		}

		// draw the screen
		Window.SwapBuffers()
		glfw.PollEvents()
	}
}
