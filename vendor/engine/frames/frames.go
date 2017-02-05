package frames

import "time"

//Frame contains present time, previus frame time, time a second ago, and values dt, and fps.
type Frame struct {
	timePrev time.Time
	timeNow  time.Time

	dtPrev float32
	dt     float32

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

var fps60 float32 = 1.0 / 80.0

//Next calculate next frame
func (f *Frame) Next() (float32, int) {

	f.timePrev = f.timeNow
	f.timeNow = time.Now()
	f.dt = float32(f.timeNow.Sub(f.timePrev).Seconds())
	// log.Println(f.dt, fps60)
	// if f.dt < fps60 {
	// 	pause := time.Duration((fps60 - f.dt) * 3500000000)
	// 	log.Println("SLEEP", pause)
	// 	time.Sleep(pause)
	// 	f.dt = fps60
	// 	// log.Println("WAKE UP")
	// }

	if f.timeFPS.Before(f.timeNow) {
		f.fps = f.prefps
		// fmt.Println("FPS:", f.fps, " DT:", f.dt)

		f.prefps = 0
		f.timeFPS = f.timeFPS.Add(time.Second)
	}

	// log.Println(f.timeNow.After(f.timeFPS))
	f.prefps++

	// dt := (f.dt + f.dtPrev) / 2
	// f.dtPrev = f.dt

	return f.dt, f.fps
}
