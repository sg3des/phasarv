package engine

import (
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

var Camera = &camera{}

type camera struct {
	mainCamera *fizzle.YawPitchCamera
}

//NewCamera create camera set it how main camera and return it
func NewCamera(eyePos mgl32.Vec3) *fizzle.YawPitchCamera {
	Camera.mainCamera = fizzle.NewYawPitchCamera(eyePos)
	render.SetCamera(Camera.mainCamera)
	return Camera.mainCamera
}

func (c *camera) SetPosition(x, y float32) {
	c.mainCamera.SetPosition(x, y, c.mainCamera.GetPosition().Z())
}

func (c *camera) GetPosition() mgl32.Vec3 {
	return c.mainCamera.GetPosition()
}
