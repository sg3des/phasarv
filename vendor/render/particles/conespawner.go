// Copyright 2016, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package particles

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
	fizzle "github.com/tbogdala/fizzle"
	renderer "github.com/tbogdala/fizzle/renderer"
)

// ConeSpawner is a particle spawner that creates particles within the
// volume of a cone as specified by the settings in the struct.
type ConeSpawner struct {
	BottomRadius float32
	TopRadius    float32
	Length       float32
	Owner        *Emitter

	volumeRenderable *fizzle.Renderable
}

// NewConeSpawner creates a new cone shaped particle spawner.
func NewConeSpawner(owner *Emitter, br, tr, cl float32) *ConeSpawner {
	cone := new(ConeSpawner)
	cone.BottomRadius = br
	cone.TopRadius = tr
	cone.Length = cl
	cone.Owner = owner
	return cone
}

// GetName returns a user friendly name for the spawner
func (cone *ConeSpawner) GetName() string {
	return "Cone Spawner"
}

// SetOwner sets the owning emitter for the spawner
func (cone *ConeSpawner) SetOwner(e *Emitter) {
	cone.Owner = e
}

// GetLocation returns the location in world space for the cone spawner.
func (cone *ConeSpawner) GetLocation() mgl.Vec3 {
	return cone.Owner.GetLocation()
}

// NewParticle creates a new particle that fits within the volume of a cone section.
func (cone *ConeSpawner) NewParticle() (p Particle) {
	// get the standard properties from the emitter itself
	p.StartTime = cone.Owner.Owner.runtime
	p.Size = cone.Owner.Properties.Size
	p.Speed = cone.Owner.Properties.Speed
	p.Color = cone.Owner.Properties.Color
	p.Acceleration = cone.Owner.Properties.Acceleration
	p.EndTime = cone.Owner.Properties.TTL + p.StartTime

	// get a random point within the bottom circle
	var bottom mgl.Vec3
	bangle := cone.Owner.rng.Float32() * math.Pi * 2.0
	bradius := cone.Owner.rng.Float32() * cone.BottomRadius
	bottom[0] = bradius * float32(math.Cos(float64(bangle)))
	bottom[2] = bradius * float32(math.Sin(float64(bangle)))

	// caculate the ratio of top to bottom size avoiding divbyzero
	var btRatio float32
	if cone.BottomRadius != 0.0 {
		btRatio = cone.TopRadius / cone.BottomRadius
	} else {
		btRatio = cone.TopRadius
	}

	// calculate the top point within the top circle
	var top mgl.Vec3
	top[0] = btRatio * bottom[0]
	top[1] = bottom[1] + cone.Length
	top[2] = btRatio * bottom[2]

	p.Velocity = top.Sub(bottom).Normalize()
	p.Velocity = cone.Owner.Properties.Rotation.Rotate(p.Velocity)

	p.Location = cone.GetLocation()
	// p.Location = cone.Owner.Properties.Rotation.Rotate(bottom)

	return p
}

// CreateRenderable creates a cached renderable for the spawner that represents
// the spawning volume for particles.
func (cone *ConeSpawner) CreateRenderable() *fizzle.Renderable {
	const circleSegments = 16
	const sideSegments = 8

	cone.volumeRenderable = fizzle.CreateWireframeConeSegmentXZ(0, 0, 0, cone.BottomRadius, cone.TopRadius, cone.Length, circleSegments, sideSegments)
	return cone.volumeRenderable
}

// DrawSpawnVolume renders a visual representation of the particle spawning volume.
func (cone *ConeSpawner) DrawSpawnVolume(r renderer.Renderer, shader *fizzle.RenderShader, projection mgl.Mat4, view mgl.Mat4, camera fizzle.Camera) {
	if cone.volumeRenderable == nil {
		cone.CreateRenderable()
	}

	// sync the position and rotation
	cone.volumeRenderable.Location = cone.Owner.Properties.Origin
	cone.volumeRenderable.LocalRotation = cone.Owner.Properties.Rotation

	r.DrawLines(cone.volumeRenderable, shader, nil, projection, view, camera)
}
