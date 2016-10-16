package engine

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

//Render object
func (o *Object) Render(perspective, view mgl32.Mat4, cam fizzle.Camera) {
	e.render.DrawRenderable(o.Node, nil, perspective, view, cam)

	if art, ok := o.GetArt("renderShape"); ok {
		art.Art.Scale = mgl32.Vec3{1, 1 - o.ShapeWidthPercent(), 1}
	}

	o.RenderChilds(perspective, view, cam)
}

//RenderChilds function for rendering childs and 2d graphics of objects
func (o *Object) RenderChilds(perspective, view mgl32.Mat4, cam fizzle.Camera) {
	for _, child := range o.ArtStatic {
		child.Art.Location = o.PositionVec3().Add(child.LocalPosition)
		if child.Line {
			e.render.DrawLines(child.Art, child.Art.Material.Shader, nil, perspective, view, cam)
		} else {
			e.render.DrawRenderable(child.Art, nil, perspective, view, cam)
		}
	}

	for _, child := range o.ArtRotate {
		xF, yF := o.VectorForward(child.LocalPosition.X())
		xS, yS := o.VectorSide(child.LocalPosition.Y(), -1.5704)
		child.Art.Location = o.PositionVec3().Add(mgl32.Vec3{xF, yF, child.LocalPosition.Z()}).Add(mgl32.Vec3{xS, yS})
		child.Art.LocalRotation = mgl32.AnglesToQuat(0, 0, o.Rotation(), 1)

		if child.Line {
			e.render.DrawLines(child.Art, child.Art.Material.Shader, nil, perspective, view, cam)
		} else {
			e.render.DrawRenderable(child.Art, nil, perspective, view, cam)
		}
	}
}