package render

import (
	"github.com/tbogdala/eweygewey"
	"github.com/tbogdala/fizzle/graphicsprovider"
)

func uiConstruct(dt float64) {
	gfx.Disable(graphicsprovider.DEPTH_TEST)
	gfx.Enable(graphicsprovider.SCISSOR_TEST)

	ui.Construct(dt)
	ui.Draw()

	gfx.Disable(graphicsprovider.SCISSOR_TEST)
	gfx.Enable(graphicsprovider.DEPTH_TEST)
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

var UI struct {
	PhysFrame   *UITextWindow
	RenderFrame *UITextWindow
}

func InitializeSystemUI() {
	UI.PhysFrame = NewTextWindow("phys", 0.79, 0.99, 0.2, 0)
	UI.RenderFrame = NewTextWindow("render", 0.79, 0.91, 0.2, 0)
}
