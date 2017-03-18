package ui

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/tbogdala/eweygewey"
	"github.com/tbogdala/eweygewey/glfwinput"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

var (
	ui *eweygewey.Manager

	PhysFrame   *UITextWindow
	RenderFrame *UITextWindow
)

func Init(gfx graphicsprovider.GraphicsProvider, window *glfw.Window) error {
	fontScale := 14
	fontFilepath := "assets/fonts/Roboto-Bold.ttf"
	fontGlyphs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890., :[]{}\\|<>;\"'~`?/-+_=()*&^%$#@!"

	// create and initialize the gui Manager
	ui = eweygewey.NewManager(gfx)
	if err := ui.Initialize(eweygewey.VertShader330, eweygewey.FragShader330, 1024, 768, 768); err != nil {
		return fmt.Errorf("Failed to initialize the user interface! reason %s", err)
	}

	glfwinput.SetInputHandlers(ui, window)

	// load a font
	_, err := ui.NewFont("Default", fontFilepath, fontScale, fontGlyphs)
	if err != nil {
		return fmt.Errorf("Failed to load the font file! reason: %s", err)
	}

	InitializeSystemUI()

	return nil
}

func Draw(dt float64) {
	ui.Construct(dt)
	ui.Draw()
}

type UITextWindow struct {
	Title string
	Text  []string
	wnd   *eweygewey.Window
}

func (w *UITextWindow) Update(wnd *eweygewey.Window) {
	wnd.Text(w.Title)
	for _, line := range w.Text {
		wnd.StartRow()
		wnd.Text(line)
	}
}

func NewTextWindow(title string, x, y, w, h float32) *UITextWindow {
	uitw := &UITextWindow{Title: title}

	uitw.wnd = ui.NewWindow(title, x, y, w, h, uitw.Update)
	uitw.wnd.ShowTitleBar = false
	uitw.wnd.IsMoveable = false
	uitw.wnd.AutoAdjustHeight = true
	return uitw
}

func InitializeSystemUI() {
	PhysFrame = NewTextWindow("phys", 0.79, 0.99, 0.2, 0)
	RenderFrame = NewTextWindow("render", 0.79, 0.91, 0.2, 0)
}
