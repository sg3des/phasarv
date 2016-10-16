package engine

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/particles"
)

var (
	textureFilepath = "assets/textures/explosion_D.png"
	particleSystem  *particles.System
)

func NewParticles() {

	// log.Println(cone.GetLocation())
	// cone.GetLocation()

	// load the particle shader
	particleShader, err := fizzle.LoadShaderProgram(particles.VertShader330, particles.FragShader330, nil)
	if err != nil {
		panic("Failed to compile and link the particle shader program! " + err.Error())
	}
	// defer particleShader.Destroy()

	particleSystem = particles.NewSystem(e.gfx)
	emitter := particleSystem.NewEmitter(nil)
	cone := particles.NewConeSpawner(emitter, 0.5, 1, 1)
	emitter.Properties.TextureFilepath = textureFilepath
	emitter.Properties.MaxParticles = 300
	emitter.Properties.SpawnRate = 40
	emitter.Properties.Size = 32.0
	emitter.Properties.Color = mgl32.Vec4{0.0, 0.9, 0.0, 1.0}
	emitter.Properties.Velocity = mgl32.Vec3{0, 1, 0}
	emitter.Properties.Acceleration = mgl32.Vec3{0, -0.1, 0}
	emitter.Properties.TTL = 3.0
	emitter.Shader = particleShader.Prog

	// load the texture
	err = emitter.LoadTexture()
	if err != nil {
		panic(err.Error())
	}

	// reset the spawner to the first known spawner instance
	emitter.Spawner = cone
	emitter.Spawner.SetOwner(emitter)
}

func loopRenderParticles(dt float32) {
	if particleSystem != nil {
		particleSystem.Update(float64(dt))
	}
}
