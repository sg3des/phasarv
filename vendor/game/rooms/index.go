package rooms

import (
	"engine"
	"scenes"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sg3des/fizzgui"
)

var sceneIndex *scenes.Scene

func Init() {
	sceneIndex = scenes.Load("room_index")
}

func Index() {
	page := new(authorize)

	c := fizzgui.NewContainer("authorize", "50px", "30%", "300px", "320px")
	c.AutoAdjustHeight = true
	page.container = c

	c.NewText("Username:")
	page.wgtUsername = c.NewInput("username", &page.user, nil)

	c.NewText("Password:")
	page.wgtPassword = c.NewInput("password", &page.pass, page.Connect)

	c.NewRow().Layout.SetHeight("10px")
	c.NewButton("Connect", page.Connect).Layout.SetWidth("100%")

	c.NewRow().Layout.SetHeight("30px")
	c.NewButton("Exit", func(wgt *fizzgui.Widget) { engine.Close() }).Layout.SetWidth("100%")
}

type authorize struct {
	container   *fizzgui.Container
	wgtUsername *fizzgui.Widget
	wgtPassword *fizzgui.Widget
	user        string
	pass        string
}

func (page *authorize) Connect(wgt *fizzgui.Widget) {
	if page.user == "" {
		page.wgtUsername.Style.BorderWidth = 2
		page.wgtUsername.Style.BorderColor = mgl32.Vec4{1, 0.3, 0.3, 1}
		return
	} else {
		page.wgtUsername.Style.BorderWidth = 0
	}

	if page.pass == "" {
		page.wgtPassword.Style.BorderWidth = 2
		page.wgtPassword.Style.BorderColor = mgl32.Vec4{1, 0.3, 0.3, 1}
		return
	} else {
		page.wgtPassword.Style.BorderWidth = 0
	}

	page.container.Close()
	Hangar(page.user)
}
