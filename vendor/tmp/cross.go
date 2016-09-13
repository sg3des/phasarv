package main

import "log"

func main() {
	log.Println(Intersection(0, 0, 10, 10, 0, 10, 10, 0))
	log.Println(Intersection(0, 0, 10, 0, 0, 10, 9, 1))
}

//*i_x, *i_y
func Intersection(p0_x, p0_y, p1_x, p1_y, p2_x, p2_y, p3_x, p3_y float32) bool {

	var s02_x, s02_y, s10_x, s10_y, s32_x, s32_y, s_numer, t_numer, denom, t float32

	s10_x = p1_x - p0_x
	s10_y = p1_y - p0_y
	s32_x = p3_x - p2_x
	s32_y = p3_y - p2_y

	denom = s10_x*s32_y - s32_x*s10_y
	if denom == 0 {
		log.Println("denom == 0  //Collinear")
		return false // Collinear
	}
	denomPositive := denom > 0

	s02_x = p0_x - p2_x
	s02_y = p0_y - p2_y
	s_numer = s10_x*s02_y - s10_y*s02_x
	if (s_numer < 0) == denomPositive {
		log.Println("s_numer < 0 ", s_numer)
		return false // No collision}
	}
	t_numer = s32_x*s02_y - s32_y*s02_x
	if (t_numer < 0) == denomPositive {
		log.Println("t_numer < 0 ", t_numer)
		return false // No collision
	}

	log.Println(denom, s_numer, t_numer, denomPositive)
	if ((s_numer >= denom) == denomPositive) || ((t_numer >= denom) == denomPositive) {
		log.Println("last")
		return false // No collision
	}
	// Collision detected
	t = t_numer / denom
	log.Println(t)
	// if (i_x != nil) {
	//     *i_x = p0_x + (t * s10_x);
	//     }
	// if (i_y != nil) {
	//     *i_y = p0_y + (t * s10_y);
	//     }

	return true
}
