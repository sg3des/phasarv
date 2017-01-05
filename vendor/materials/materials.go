package materials

import (
	"assets"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

var (
	Materials = make(map[string]*fizzle.Material)
)

type Instruction struct {
	Name            string
	Shader, Texture string
	DiffColor       mgl32.Vec4
	SpecLevel       float32
}

func (i *Instruction) Create() *fizzle.Material {
	// m, ok := Materials[i.Name]
	// if ok {
	// 	return m
	// }

	if i.Shader == "" {
		i.Shader = "color"
	}
	if i.Texture == "" {
		i.Texture = "gray"
	}

	m := fizzle.NewMaterial()

	m.Shader = assets.GetShader(i.Shader)

	m.DiffuseTex = assets.GetTexture(i.Texture).Diffuse
	m.NormalsTex = assets.GetTexture(i.Texture).Normals

	if i.DiffColor.Len() != 0 {
		m.DiffuseColor = i.DiffColor
	}

	m.SpecularColor = mgl32.Vec4{i.SpecLevel, i.SpecLevel, i.SpecLevel, 1}
	m.Shininess = i.SpecLevel

	Materials[i.Name] = m

	return m
}
