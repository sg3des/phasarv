package phys

const (
	errorBias = 0.00179701029991443 //cpfpow(1.0f - 0.1f, 60.0f)
)

type ConstraintCallback interface {
	CollisionPreSolve(constraint Constraint)
	CollisionPostSolve(constraint Constraint)
}

type Constraint interface {
	Constraint() *BasicConstraint
	PreSolve()
	PostSolve()
	PreStep(dt float32)
	ApplyCachedImpulse(dt_coef float32)
	ApplyImpulse()
	Impulse() float32
}

type BasicConstraint struct {
	BodyA, BodyB    *Body
	space           *Space
	MaxForce        float32
	ErrorBias       float32
	MaxBias         float32
	CallbackHandler ConstraintCallback
	UserData        Data
}

func NewConstraint(a, b *Body) BasicConstraint {
	return BasicConstraint{BodyA: a, BodyB: b, MaxForce: Inf, MaxBias: Inf, ErrorBias: errorBias}
}

func (this *BasicConstraint) Constraint() *BasicConstraint {
	return this
}

func (this *BasicConstraint) PreStep(dt float32) {
	panic("empty constraint")
}

func (this *BasicConstraint) ApplyCachedImpulse(dt_coef float32) {
	panic("empty constraint")
}

func (this *BasicConstraint) ApplyImpulse() {
	panic("empty constraint")
}

func (this *BasicConstraint) Impulse() float32 {
	panic("empty constraint")
}

func (this *BasicConstraint) PreSolve() {
	if this.CallbackHandler != nil {
		this.CallbackHandler.CollisionPreSolve(this)
	}
}

func (this *BasicConstraint) PostSolve() {
	if this.CallbackHandler != nil {
		this.CallbackHandler.CollisionPostSolve(this)
	}
}
