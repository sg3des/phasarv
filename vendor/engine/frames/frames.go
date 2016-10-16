package frames

import "time"

//Frame contains present time, previus frame time, time a second ago, and values dt, and fps.
type Frame struct {
	timePrev time.Time
	timeNow  time.Time

	dt float32

	timeFPS time.Time
	prefps  int
	fps     int
}

//NewFrame prepare and return frame struct
func NewFrame() *Frame {
	f := &Frame{
		timeNow: time.Now(),
		timeFPS: time.Now(),
	}
	return f
}

//Next calculate next frame
func (f *Frame) Next() (float32, int) {

	f.timePrev = f.timeNow
	f.timeNow = time.Now()
	f.dt = float32(f.timeNow.Sub(f.timePrev).Seconds())

	if f.timeFPS.Before(f.timeNow) {
		f.fps = f.prefps
		// fmt.Println("FPS:", f.fps, " DT:", f.dt)

		f.prefps = 0
		f.timeFPS = f.timeFPS.Add(time.Second)
	}

	// log.Println(f.timeNow.After(f.timeFPS))
	f.prefps++

	return f.dt, f.fps
}
