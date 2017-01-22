package game

import "phys/vect"

//Render flag if it false, graphics elements(bars,aims,trails,etc...) should not be initialized.
var Render bool

//NetPacket structure of standard network packet
type NetPacket struct {
	Vel  vect.Vect
	AVel float32
	Pos  vect.Vect
}
