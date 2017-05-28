package phys

import (
	"phys/transform"
	"phys/vect"
)

type ShapeType int

const (
	ShapeType_Circle = iota
	ShapeType_Segment
	ShapeType_Polygon
	ShapeType_Box
	numShapes
)

func (st ShapeType) ToString() string {
	switch st {
	case ShapeType_Circle:
		return "Circle"
	case ShapeType_Segment:
		return "Segment"
	case ShapeType_Polygon:
		return "Polygon"
	case ShapeType_Box:
		return "Box"
	default:
		return "Unknown"
	}
	panic("never reached")
}

type ShapeClass interface {
	ShapeType() ShapeType
	// Update the shape with the new transform and compute the AABB.
	update(xf transform.Transform) AABB
	// Returns if the given point is located inside the shape.
	TestPoint(point vect.Vect) bool

	Moment(mass float32) float32

	Clone(s *Shape) ShapeClass
	//marshalShape(shape *Shape) ([]byte, error)
	//unmarshalShape(shape *Shape, data []byte) error
}

// Returns shape.ShapeClass as CircleShape or nil.
func (shape *Shape) GetAsCircle() *CircleShape {
	if circle, ok := shape.ShapeClass.(*CircleShape); ok {
		return circle
	}

	return nil
}

// Returns shape.ShapeClass as PolygonShape or nil.
func (shape *Shape) GetAsPolygon() *PolygonShape {
	if poly, ok := shape.ShapeClass.(*PolygonShape); ok {
		return poly
	}

	return nil
}

// Returns shape.ShapeClass as SegmentShape or nil.
func (shape *Shape) GetAsSegment() *SegmentShape {
	if seg, ok := shape.ShapeClass.(*SegmentShape); ok {
		return seg
	}

	return nil
}

// Returns shape.ShapeClass as BoxShape or nil.
func (shape *Shape) GetAsBox() *BoxShape {
	if box, ok := shape.ShapeClass.(*BoxShape); ok {
		return box
	}

	return nil
}
