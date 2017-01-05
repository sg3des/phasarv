package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle/renderer/forward"
)

type Light struct {
	Pos mgl32.Vec3
	Dir mgl32.Vec3

	Strength float32
	Specular float32

	ShadowSize int32
	Direct     bool

	LightNode *forward.Light
}

func (l *Light) Create() {
	if l.Dir.Len() == 0 {
		l.Dir = l.Pos.Mul(-1)
	}

	if l.Direct {
		l.LightNode = render.NewDirectionalLight(l.Dir)
	} else {
		l.LightNode = render.NewPointLight(l.Pos)
	}
	if l.ShadowSize > 0 {
		l.LightNode.CreateShadowMap(l.ShadowSize, 1, 400, l.Dir)
	}

	l.LightNode.Direction = l.Dir
	l.LightNode.Position = l.Pos

	l.LightNode.DiffuseColor = mgl32.Vec4{0.9, 0.9, 0.9, 1}
	l.LightNode.DiffuseIntensity = l.Strength
	l.LightNode.Strength = l.Strength
	l.LightNode.SpecularIntensity = l.Specular

	render.ActiveLights[render.GetActiveLightCount()] = l.LightNode
}
