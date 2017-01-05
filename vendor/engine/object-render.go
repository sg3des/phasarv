package engine

// //Render object
// func (o *Object) Render(perspective, view mgl32.Mat4, cam fizzle.Camera) {
// 	// log.Println(o.Name)
// 	e.render.DrawRenderable(o.Node, nil, perspective, view, cam)

// 	if art, ok := o.GetArt("renderShape"); ok {
// 		art.Node.Scale = mgl32.Vec3{1, 1 - o.ShapeWidthPercent(), 1}
// 	}

// 	o.Renderarts(perspective, view, cam)
// }

// //Renderarts function for rendering arts and 2d graphics of objects
// func (o *Object) Renderarts(perspective, view mgl32.Mat4, cam fizzle.Camera) {
// 	for _, art := range o.ArtStatic {
// 		if art.Node != nil {
// 			art.Node.Location = o.PositionVec3().Add(art.Param.Pos.Vec3())
// 			if art.RenderLine {
// 				e.render.DrawLines(art.Node, art.Node.Material.Shader, nil, perspective, view, cam)
// 			} else {
// 				// log.Println(art.Name, art.Node)
// 				// log.Println(art.Node.Material)
// 				e.render.DrawRenderable(art.Node, nil, perspective, view, cam)
// 			}
// 		}
// 	}

// 	for _, art := range o.ArtRotate {
// 		if art.Node != nil {
// 			xF, yF := o.VectorForward(art.Param.Pos.X)
// 			xS, yS := o.VectorSide(art.Param.Pos.Y, -1.5704)
// 			art.Node.Location = o.PositionVec3().Add(mgl32.Vec3{xF, yF, art.Param.Pos.Z}).Add(mgl32.Vec3{xS, yS})
// 			art.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, o.Rotation(), 1)

// 			if art.RenderLine {
// 				e.render.DrawLines(art.Node, art.Node.Material.Shader, nil, perspective, view, cam)
// 			} else {
// 				e.render.DrawRenderable(art.Node, nil, perspective, view, cam)
// 			}
// 		}
// 	}
// }
