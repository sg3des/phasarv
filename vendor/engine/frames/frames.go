package frames

import (
	"fmt"
	"time"
)

//Frame contains present time, previus frame time, time a second ago, and values dt, and fps.
type Frame struct {
	timePrev time.Time
	timeNow  time.Time

	dt float32

	timeFPS time.Time
	fps     int
}

//NewFrame prepare and return frame struct
func NewFrame() *Frame {
	f := &Frame{
		timePrev: time.Now(),
		timeFPS:  time.Now(),
	}
	return f
}

//Next calculate next frame
func (f *Frame) Next() float32 {

	f.timePrev = f.timeNow
	f.timeNow = time.Now()
	f.dt = float32(f.timeNow.Sub(f.timePrev).Seconds())

	// log.Println(f.timeNow.After(f.timeFPS))
	f.fps++

	return f.dt
}

//FPS return fps value
func (f *Frame) FPS() {
	if f.timeFPS.Before(f.timeNow) {
		fmt.Println("FPS:", f.fps, " DT:", f.dt)
		f.fps = 0
		f.timeFPS = f.timeFPS.Add(time.Second)
	}
}

//DT return delta time
func (f *Frame) DT() float32 {
	return f.dt
}
