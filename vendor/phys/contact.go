package phys

import (
	"phys/vect"
)

type Contact struct {
	p, n vect.Vect
	dist float32

	r1, r2               vect.Vect
	nMass, tMass, bounce float32

	jnAcc, jtAcc, jBias float32
	bias                float32

	hash HashValue
}

func (con *Contact) reset(pos, norm vect.Vect, dist float32, hash HashValue) {
	con.p = pos
	con.n = norm
	con.dist = dist
	con.hash = hash

	con.jnAcc = 0.0
	con.jtAcc = 0.0
	con.jBias = 0.0
}

func (con *Contact) Normal() vect.Vect {
	return con.n
}

func (con *Contact) Position() vect.Vect {
	return con.p
}
