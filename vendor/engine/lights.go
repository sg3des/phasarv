package engine

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle/renderer/forward"
)

var (
	lights = make(map[*forward.Light]bool)
)

//NewLight Create new list
func NewLight(diffuse, attenuation float32, shadowSize int32) *forward.Light {
	light := e.render.NewPointLight(mgl32.Vec3{0, 0, 10})
	light.DiffuseColor = mgl32.Vec4{0.9, 0.9, 0.9, 1.0}
	light.DiffuseIntensity = diffuse
	light.AmbientIntensity = 0.5
	light.Strength = diffuse
	light.DiffuseIntensity = 1

	light.LinearAttenuation = attenuation

	e.render.ActiveLights[len(lights)] = light

	// if shadow {
	light.CreateShadowMap(shadowSize, 0.1, 100.0, mgl32.Vec3{-1, -1, -10})
	// }

	lights[light] = true

	return light
}

func NewSun() *forward.Light {
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
