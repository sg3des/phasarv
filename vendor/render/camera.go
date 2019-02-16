package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

// var Camera = &camera{}

// type camera struct {
// 	mainCamera *fizzle.YawPitchCamera
// }

//NewCamera create camera set it how main camera and return it
func NewCamera(eyePos mgl32.Vec3) *fizzle.YawPitchCamera {
	camera = fizzle.NewYawPitchCamera(eyePos)
	// SetCamera(camera)
	return camera
}

func SetCameraPosition(x, y float32) {
	camera.SetPosition(x, y, camera.GetPosition().Z())

	SetLightPositions(x, y)
}

func SetLightPositions(x, y float32) {
	count := render.GetActiveLightCount()
	for i := 0; i < count; i++ {
		light := render.ActiveLights[i]

		light.Position = mgl32.Vec3{x, y, light.Position[2]}
	}

}

func GetCameraPosition() mgl32.Vec3 {
	return camera.GetPosition()
}

func GetCameraViewMatrix() mgl32.Mat4 {
	return camera.GetViewMatrix()
}
