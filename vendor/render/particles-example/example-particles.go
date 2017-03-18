package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	glfw "github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	gui "github.com/tbogdala/eweygewey"

	fizzle "github.com/tbogdala/fizzle"
	graphics "github.com/tbogdala/fizzle/graphicsprovider"
	opengl "github.com/tbogdala/fizzle/graphicsprovider/opengl"
	particles "github.com/tbogdala/fizzle/particles"
	forward "github.com/tbogdala/fizzle/renderer/forward"
)

var (
	windowWidth     = 1280
	windowHeight    = 720
	mainWindow      *glfw.Window
	renderer        *forward.ForwardRenderer
	textureFilepath = "explosion00.png"
)

// spawnerPrototypes keeps track of possible spawner interface implementations
// to switch between
type spawnerPrototypes struct {
	Name string
	particles.ParticleSpawner
	RenderUI func(wnd *gui.Window)
}

var (
	// created instances of particle spawners
	knownSpawners spawnerPrototypes
)

// GLFW event handling must run on the main OS thread. If this doesn't get
// locked down, you will likely see random crashes on memory access while
// running the application after a few seconds.
//
// So on initialization of the module, lock the OS thread for this goroutine.
func init() {
	runtime.LockOSThread()
}

// initSpawners create prototype instances of all known spawner types
func initSpawners() {
	// var knownSpawners = spawnerPrototypes

	cone := particles.NewConeSpawner(nil, 0.5, 1, 1)
	knownSpawners = spawnerPrototypes{Name: cone.GetName(), ParticleSpawner: cone, RenderUI: func(wnd *gui.Window) {
		const textWidth = 0.33
		wnd.RequestItemWidthMin(textWidth)
		wnd.Text("Top Radius")
		wnd.DragSliderUFloat("tradius", 0.1, &cone.TopRadius)

		wnd.StartRow()
		wnd.RequestItemWidthMin(textWidth)
		wnd.Text("Bottom Radius")
		wnd.DragSliderUFloat("bradius", 0.1, &cone.BottomRadius)

		wnd.StartRow()
		wnd.RequestItemWidthMin(textWidth)
		wnd.Text("Length")
		wnd.DragSliderUFloat("conelength", 0.1, &cone.Length)
	}}
}

func main() {
	// start off by initializing the GL and GLFW libraries and creating a window.
	w, gfx := initGraphics("Particle Editor", windowWidth, windowHeight)
	mainWindow = w

	// create prototype instances of all known spawner types
	initSpawners()

	/////////////////////////////////////////////////////////////////////////////
	// make a window that will render the particle system
	const particleWindowSize = 512

	renderer = forward.NewForwardRenderer(gfx)
	renderer.ChangeResolution(particleWindowSize, particleWindowSize)
	defer renderer.Destroy()

	// load the particle shader
	particleShader, err := fizzle.LoadShaderProgram(particles.VertShader330, particles.FragShader330, nil)
	if err != nil {
		panic("Failed to compile and link the particle shader program! " + err.Error())
	}
	defer particleShader.Destroy()

	// load the color shader
	colorShader, err := forward.CreateColorShader()
	if err != nil {
		panic("Failed to compile and link the color shader program! " + err.Error())
	}
	defer colorShader.Destroy()

	// create a particle system
	particleSystem := particles.NewSystem(gfx)
	emitter := particleSystem.NewEmitter(nil)
	emitter.Properties.TextureFilepath = textureFilepath
	emitter.Properties.MaxParticles = 300
	emitter.Properties.SpawnRate = 40
	emitter.Properties.Size = 32.0
	emitter.Properties.Color = mgl.Vec4{0.0, 0.9, 0.0, 1.0}
	emitter.Properties.Velocity = mgl.Vec3{0, 1, 0}
	emitter.Properties.Acceleration = mgl.Vec3{0, -0.1, 0}
	emitter.Properties.TTL = 3.0
	emitter.Shader = particleShader.Prog

	// load the texture
	err = emitter.LoadTexture()
	if err != nil {
		panic(err.Error())
	}

	// reset the spawner to the first known spawner instance
	emitter.Spawner = knownSpawners
	emitter.Spawner.SetOwner(emitter)

	// setup the camera to look at the cube
	camera := fizzle.NewOrbitCamera(mgl.Vec3{0, 0, 0}, math.Pi/2.0, 5.0, math.Pi/2.0)

	// /////////////////////////////////////////////////////////////////////////////
	// // loop until something told the mainWindow that it should close
	// // set some OpenGL flags
	gfx.Enable(graphics.CULL_FACE)
	gfx.Enable(graphics.DEPTH_TEST)
	gfx.Enable(graphics.PROGRAM_POINT_SIZE)
	gfx.Enable(graphics.BLEND)
	gfx.BlendFunc(graphics.SRC_ALPHA, graphics.ONE_MINUS_SRC_ALPHA)

	timeFPS := time.Now().Add(time.Second)
	lastFrame := time.Now()
	var i int
	for !mainWindow.ShouldClose() {
		i++
		// calculate the difference in time to control rotation speed
		thisFrame := time.Now()
		frameDelta := thisFrame.Sub(lastFrame).Seconds()

		if time.Now().After(timeFPS) {
			fmt.Println("FPS:", i)
			i = 0
			timeFPS = time.Now().Add(time.Second)
		}

		// update the data for the application
		particleSystem.Update(frameDelta)

		// clear the screen
		//width, height := renderer.GetResolution()
		gfx.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
		clearColor := gui.ColorIToV(114, 144, 154, 255)
		gfx.ClearColor(clearColor[0], clearColor[1], clearColor[2], clearColor[3])
		gfx.Clear(graphics.COLOR_BUFFER_BIT | graphics.DEPTH_BUFFER_BIT)

		// draw the user interface
		// uiman.Construct(frameDelta)
		// uiman.Draw()

		// rotate the cube and sphere around the Y axis at a speed of radsPerSec
		// gfx.ClearColor(0.0, 0.0, 0.0, 1.0)
		// gfx.Clear(graphics.COLOR_BUFFER_BIT | graphics.DEPTH_BUFFER_BIT)

		perspective := mgl.Perspective(mgl.DegToRad(60.0), float32(particleWindowSize)/float32(particleWindowSize), 0.1, 50.0)
		view := camera.GetViewMatrix()
		particleSystem.Draw(perspective, view)

		// draw the emitter volumes
		for _, e := range particleSystem.Emitters {
			e.Spawner.CreateRenderable()
			e.Spawner.DrawSpawnVolume(renderer, colorShader, perspective, view, camera)
		}

		// draw the screen
		mainWindow.SwapBuffers()

		// advise GLFW to poll for input. without this the window appears to hang.
		glfw.PollEvents()

		// update our last frame time
		lastFrame = thisFrame
	}
}

// initGraphics creates an OpenGL window and initializes the required graphics libraries.
// It will either succeed or panic.
func initGraphics(title string, w int, h int) (*glfw.Window, graphics.GraphicsProvider) {
	// GLFW must be initialized before it's called
	err := glfw.Init()
	if err != nil {
		panic("Can't init glfw! " + err.Error())
	}

	// request a OpenGL 3.3 core context
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// do the actual window creation
	mainWindow, err = glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		panic("Failed to create the main window! " + err.Error())
	}
	// mainWindow.SetSizeCallback(onWindowResize)
	mainWindow.MakeContextCurrent()

	// disable v-sync for max draw rate
	glfw.SwapInterval(0)

	// initialize OpenGL
	gfx, err := opengl.InitOpenGL()
	if err != nil {
		panic("Failed to initialize OpenGL! " + err.Error())
	}
	fizzle.SetGraphics(gfx)

	return mainWindow, gfx
}

// onWindowResize is called when the window changes size
// func onWindowResize(w *glfw.Window, width int, height int) {
// 	uiman.AdviseResolution(int32(width), int32(height))
// 	//renderer.ChangeResolution(int32(width), int32(height))
// }
