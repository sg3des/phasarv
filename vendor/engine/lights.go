package engine

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle/renderer/forward"
)

var (
	lights     = make(map[*forward.Light]bool)
	cNewLights = make(chan ParamLight)
)

type ParamLight struct {
	Sun bool

	Pos mgl32.Vec3

	Strength float32
	Specular float32

	ShadowSize int32

	c chan *forward.Light
}

func NewLight(p ParamLight) (light *forward.Light) {
	p.c = make(chan *forward.Light)
	cNewLights <- p

	return <-p.c
}

func (p *ParamLight) create() {
	var light *forward.Light

	dir := p.Pos.Mul(-1)
	if p.Sun {
		light = e.render.NewDirectionalLight(dir)
	} else {
		light = e.render.NewPointLight(p.Pos)
	}

	light.Direction = dir
	light.Position = p.Pos

	// light.Position = p.Pos
	// light.Direction = p.Dir

	// light.Direction = p.Pos
	// light.Position = p.Dir
	// light.LinearAttenuation = p.Diffuse

	light.DiffuseColor = mgl32.Vec4{0.9, 0.9, 0.9, 1}
	light.DiffuseIntensity = p.Strength
	light.Strength = p.Strength
	light.SpecularIntensity = p.Specular
	// light.DiffuseIntensity = p.Diffuse

	light.CreateShadowMap(p.ShadowSize, 0.5, 400, dir)

	e.render.ActiveLights[len(lights)] = light

	lights[light] = true

	p.c <- light
}

// //NewLight Create new list
func createLight() *forward.Light {
	light := e.render.NewPointLight(mgl32.Vec3{0, 0, 10})
	light.DiffuseColor = mgl32.Vec4{0.9, 0.9, 0.9, 1.0}
	// light.DiffuseIntensity = 10
	light.AmbientIntensity = 0.5
	light.Strength = 10
	light.DiffuseIntensity = 1

	light.LinearAttenuation = 1

	e.render.ActiveLights[len(lights)] = light

	// if shadow {
	light.CreateShadowMap(2, 0.1, 100.0, mgl32.Vec3{-1, -1, -10})
	// }

	lights[light] = true

	return light
}

func createSun() *forward.Light {
	pos := mgl32.Vec3{-30, 30, 100}
	light := e.render.NewDirectionalLight(pos)
	light.DiffuseColor = mgl32.Vec4{1, 1, 1, 1}
	light.Direction = mgl32.Vec3{30, -30, -100}
	light.Strength = 0.5
	// light.DiffuseIntensity = 0.5
	light.SpecularIntensity = 0.5
	// light.AmbientIntensity = 0.5
	// light.LinearAttenuation = 1.0

	light.Position = pos

	light.CreateShadowMap(8192, 1, 400.0, mgl32.Vec3{30, -30, -100})

	e.render.ActiveLights[len(lights)] = light

	lights[light] = true

	return light
}
