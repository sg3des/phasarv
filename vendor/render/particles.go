package render

import (
	"assets"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/particles"
)

func NewTrail() {
	// load the particle shader
	particleShader, err := fizzle.LoadShaderProgram(particles.VertShader330, particles.FragShader330, nil)
	if err != nil {
		panic("Failed to compile and link the particle shader program! " + err.Error())
	}
	defer particleShader.Destroy()

	particleSystem := particles.NewSystem(gfx)
	emitter := particleSystem.NewEmitter(nil)
	emitter.Properties.TextureFilepath = assets.GetTexture("explosion00.png").Path
	emitter.Properties.MaxParticles = 300
	emitter.Properties.SpawnRate = 40
	emitter.Properties.Size = 32.0
	emitter.Properties.Color = mgl32.Vec4{0.0, 0.9, 0.0, 1.0}
	emitter.Properties.Velocity = mgl32.Vec3{0, 1, 0}
	emitter.Properties.Acceleration = mgl32.Vec3{0, -0.1, 0}
	emitter.Properties.TTL = 3.0
	emitter.Shader = particleShader.Prog
}
