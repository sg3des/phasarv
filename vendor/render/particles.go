package render

import (
	"log"
	"materials"
	"math"
	"math/rand"
	"phys/vect"
	"point"

	"github.com/chewxy/math32"
	"github.com/go-gl/mathgl/mgl32"
)

var Particles []*Particle

func drawParticles(dt float32) {
	for i := 0; i < len(Particles); i++ {
		p := Particles[i]

		if len(p.objects) == 0 {
			copy(Particles[i:], Particles[i+1:])
			Particles[len(Particles)-1] = nil // or the zero value of T
			Particles = Particles[:len(Particles)-1]
			i--
			continue
		}

		p.calculate(dt)
		p.render()
	}

	// for i, p := range Particles {
	// 	if len(p.objects) == 0 {
	// 		Particles[i] = nil
	// 		Particles = append(Particles[:i], Particles[i+1:]...)
	// 		continue
	// 	}
	// 	p.render()
	// }
}

type Particle struct {
	parent *Renderable

	i       int
	objects []*Renderable
	points  []trailPoint
	offset  mgl32.Vec3

	dist            float32
	prevParentPoint mgl32.Vec3
}

type trailPoint struct {
	pos   mgl32.Vec3
	Alpha float32
	// X, Y, Angle, Alpha float32
}

// Norm returns the L1-Norm of a Quaternion (W,X,Y,Z) -> Sqrt(W*W+X*X+Y*Y+Z*Z)
func Norm(q mgl32.Quat) float32 {
	x, y, z := q.V.Elem()
	return math32.Sqrt(q.W*q.W + x*x + y*y + z*z)
}

// Unit returns the Quaternion rescaled to unit-L1-norm
func Unit(q mgl32.Quat) mgl32.Quat {
	// x, y, z := q.V.Elem()
	k := Norm(q)
	q.W = q.W / k
	q.V[0] = q.V[0] / k
	q.V[1] = q.V[1] / k
	q.V[2] = q.V[2] / k
	return q
	// return mgl32.Quat{
	// 	W: q.W,
	// 	V: mgl32.Vec3[0]
	// }
	// return mgl32.Quat{q.W / k, x / k, y / k, z / k}
}

func Euler(q mgl32.Quat) (float32, float32, float32) {
	q = Unit(q)
	x, y, z := q.V.Elem()
	phi := math32.Atan2(2*(q.W*x+y*z), 1-2*(x*x+y*y))
	theta := math32.Asin(2 * (q.W*y - z*x))
	psi := math32.Atan2(2*(x*y+q.W*z), 1-2*(y*y+z*z))
	return phi, theta, psi
}

// func QuatAxis(q mgl32.Quat) (float32,float32,float32) {
// 	angle := 2 * math32.Acos(q.W)
// 	x = qx / sqrt(1-qw*qw)
// 	y = qy / sqrt(1-qw*qw)
// 	z = qz / sqrt(1-qw*qw)
// }

func AngleBetweenPoints(a, b mgl32.Vec3) float32 {
	return float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
}

func (p *Particle) render() {
	for _, o := range p.objects {
		o.render()
	}
}

func (p *Particle) calculate(dt float32) {
	if p == nil || p.parent == nil || p.parent.Body == nil {
		return
	}

	q := p.parent.Body.LocalRotation
	_, _, z := Euler(q)

	// q = q.Normalize()
	// angle := 2 * math32.Acos(q.W) //+ p.offset.Y()

	// log.Println(x, y, z, angle)

	v2 := vect.FromAngle(z + p.offset.Y())

	v2.Mult(p.offset.X())

	pos := p.parent.Body.Location
	pos = pos.Add(v2.Vec3())

	p.objects[0].Body.Location = pos
	p.objects[0].Body.Scale = mgl32.Vec3{1.1 + rand.Float32(), 2.2 + rand.Float32(), 1}

	dist := pos.Sub(p.prevParentPoint).Len()
	if dist > p.dist {
		p.prevParentPoint = pos

		point := trailPoint{pos, 1}
		p.points = append([]trailPoint{point}, p.points[:len(p.points)]...)

		for i, o := range p.objects {
			o.Body.Location = p.points[i].pos
			if i == 0 {
				o.Body.LocalRotation = p.parent.Body.LocalRotation
			} else {
				angle := AngleBetweenPoints(o.Body.Location, p.points[i-1].pos)
				o.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, angle, 1)
			}
		}
	}

	var summAlpha float32
	for i, o := range p.objects {
		p.points[i].Alpha = p.points[i].Alpha - 500/float32(len(p.points))*dt
		if p.points[i].Alpha < 0 {
			p.points[i].Alpha = 0
		}
		summAlpha += p.points[i].Alpha
		if o.Body == nil {
			log.Println("return - wtf?")
			return
		}

		o.Body.Material.DiffuseColor[3] = p.points[i].Alpha
	}

	p.objects[0].Body.Material.DiffuseColor[0] = 1
	p.objects[0].Body.Material.DiffuseColor[1] = 0.3
	p.objects[0].Body.Material.DiffuseColor[2] = 0

	p.objects[1].Body.Material.DiffuseColor[0] = 1
	p.objects[1].Body.Material.DiffuseColor[1] = 0.6
	p.objects[1].Body.Material.DiffuseColor[2] = 0.1
	p.objects[1].Body.Scale = mgl32.Vec3{1.1, 1.5, 1}

	// log.Println(summAlpha, p.parent == nil, p.parent.needDestroy)
	if summAlpha <= 1 && (p.parent == nil || p.parent.needDestroy) {
		p.objects = nil
	}
}

func (parent *Renderable) NewTrail(offset mgl32.Vec3, count int, size point.P) {

	p := &Particle{
		parent: parent,
		offset: offset,
		dist:   size.Y * 0.7,
	}

	for i := 0; i < count; i++ {
		plane := &Renderable{
			Transparent: true,
			P:           &point.Param{Size: size},
			RI: &Instruction{
				MeshName:    "plane",
				Material:    &materials.Instruction{Name: "laser", Texture: "laser", Shader: "colorblend", DiffColor: mgl32.Vec4{1, 1, 1, 1}},
				Transparent: true,
			},
		}
		p.objects = append(p.objects, plane)
		p.points = append(p.points, trailPoint{})
	}

	Particles = append(Particles, p)
	parent.particles = append(Particles, p)
}

func (p *Particle) Destroy() {
	if p == nil {
		return
	}
	for _, o := range p.objects {
		o.needDestroy = true
	}
	// p.objects = nil
}

// func NewTrail(pos mgl32.Vec3, count, size float32) {

// }
