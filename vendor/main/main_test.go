package main

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCursor(t *testing.T) {
	fmt.Println(getCursorPos(0, 0, 800, 600, mgl32.Vec3{0, 20, 0}))
	fmt.Println(getCursorPos(50, 100, 800, 600, mgl32.Vec3{0, 20, 0}))
	fmt.Println(getCursorPos(600, 300, 800, 600, mgl32.Vec3{0, 20, 0}))
}
