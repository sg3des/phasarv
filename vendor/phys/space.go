package phys

import (
	"errors"
	"fmt"
	"log"
	"phys/transform"
	"phys/vect"
	//"github.com/davecgh/go-spew/spew"
	"math"
	"time"
)

const ArbiterBufferSize = 1000
const ContactBufferSize = ArbiterBufferSize * MaxPoints

type Space struct {

	/// Number of iterations to use in the impulse solver to solve contacts.
	Iterations int

	/// Gravity to pass to rigid bodies when integrating velocity.
	Gravity vect.Vect

	/// Linear damping rate expressed as the fraction of linear velocity bodies retain each second.
	/// A value of 0.9 would mean that each body's velocity will drop 10% per second.
	/// The default value is 1.0, meaning no damping is applied.
	/// @note This damping value is different than those of cpDampedSpring and cpDampedRotarySpring.
	LinearDamping float32

	/// Angular damping is the same as linear damping, but for angular velocity
	AngularDamping float32

	/// Speed threshold for a body to be considered idle.
	/// The default value of 0 means to let the space guess a good threshold based on gravity.
	idleSpeedThreshold float32

	/// Time a group of bodies must remain idle in order to fall asleep.
	/// Enabling sleeping also implicitly enables the the contact graph.
	/// The default value of INFINITY disables the sleeping algorithm.
	sleepTimeThreshold float32

	/// Amount of encouraged penetration between colliding shapes.
	/// Used to reduce oscillating contacts and keep the collision cache warm.
	/// Defaults to 0.1. If you have poor simulation quality,
	/// increase this number as much as possible without allowing visible amounts of overlap.
	collisionSlop float32

	/// Determines how fast overlapping shapes are pushed apart.
	/// Expressed as a fraction of the error remaining after each second.
	/// Defaults to pow(1.0 - 0.1, 60.0) meaning that Chipmunk fixes 10% of overlap each frame at 60Hz.
	collisionBias float32

	/// Number of frames that contact information should persist.
	/// Defaults to 3. There is probably never a reason to change this value.
	collisionPersistence int64

	/// Rebuild the contact graph during each step. Must be enabled to use the cpBodyEachArbiter() function.
	/// Disabled by default for a small performance boost. Enabled implicitly when the sleeping feature is enabled.
	enableContactGraph bool

	curr_dt float32

	Constraints []Constraint

	Bodies             []*Body
	sleepingComponents []*Body
	deleteBodies       []*Body

	stamp time.Duration

	staticShapes *SpatialIndex
	activeShapes *SpatialIndex

	cachedArbiters map[HashPair]*Arbiter
	Arbiters       []*Arbiter

	ArbiterBuffer []*Arbiter
	ContactBuffer [][]*Contact

	ApplyImpulsesTime time.Duration
	ReindexQueryTime  time.Duration
	StepTime          time.Duration
}

type ContactBufferHeader struct {
	stamp       time.Duration
	next        *ContactBufferHeader
	numContacts int
}

type ContactBuffer struct {
	header   ContactBufferHeader
	contacts [256]Contact
}

func NewSpace() (space *Space) {

	space = &Space{}
	space.Iterations = 20

	space.Gravity = vect.Vector_Zero

	space.LinearDamping = 1.0
	space.AngularDamping = 1.0

	space.collisionSlop = 0.5
	space.collisionBias = float32(math.Pow(1.0-0.1, 60))
	space.collisionPersistence = 3

	space.Constraints = make([]Constraint, 0)

	space.Bodies = make([]*Body, 0)
	space.deleteBodies = make([]*Body, 0)
	space.sleepingComponents = make([]*Body, 0)

	space.staticShapes = NewBBTree(nil)
	space.activeShapes = NewBBTree(space.staticShapes)
	space.cachedArbiters = make(map[HashPair]*Arbiter)
	space.Arbiters = make([]*Arbiter, 0)
	space.ArbiterBuffer = make([]*Arbiter, ArbiterBufferSize)

	for i := 0; i < len(space.ArbiterBuffer); i++ {
		space.ArbiterBuffer[i] = newArbiter()
	}

	space.ContactBuffer = make([][]*Contact, ContactBufferSize)

	for i := 0; i < len(space.ContactBuffer); i++ {
		var contacts []*Contact = make([]*Contact, MaxPoints)

		for i := 0; i < MaxPoints; i++ {
			contacts[i] = &Contact{}
		}
		space.ContactBuffer[i] = contacts
	}
	/*
		for i := 0; i < 8; i++ {
			go space.MultiThreadTest()
		}
	*/
	return
}

func (space *Space) Destory() {
	fmt.Println("Destory is depricated, used Destroy instead.")
	space.Destroy()
}

func (space *Space) Destroy() {
	space.Bodies = nil
	space.sleepingComponents = nil
	space.staticShapes = nil
	space.activeShapes = nil
	space.cachedArbiters = nil
	space.Arbiters = nil
	space.ArbiterBuffer = nil
	space.ContactBuffer = nil
}

func (space *Space) Step(dt float32) {

	// don't step if the timestep is 0!
	if dt == 0 {
		log.Println("WARNING: dt is 0")
		return
	}

	// stepStart := time.Now()

	// bodies := space.Bodies

	for _, arb := range space.Arbiters {
		arb.state = arbiterStateNormal
	}

	space.Arbiters = space.Arbiters[0:0]

	// prev_dt := space.curr_dt
	// space.curr_dt = dt

	space.stamp++

	for _, body := range space.Bodies {
		// for _, s := range body.Shapes {
		// 	log.Println(s.Group, body.Mass())
		// }
		if body.Enabled {
			body.UpdatePosition(dt)
			body.UpdateShapes()
		}
	}

	// for _, body := range bodies {
	// 	if body.Enabled {
	// 		body.UpdatePosition(dt)
	// 	}
	// }

	// for _, body := range bodies {
	// 	if body.Enabled {
	// 		body.UpdateShapes()
	// 	}
	// }

	start := time.Now()
	space.activeShapes.ReindexQuery(func(a, b Indexable) {
		SpaceCollideShapes(a.Shape(), b.Shape(), space)
	})
	space.ReindexQueryTime = time.Since(start)

	//axc := space.activeShapes.SpatialIndexClass.(*BBTree)
	//PrintTree(axc.root)

	for h, arb := range space.cachedArbiters {
		ticks := space.stamp - arb.stamp
		deleted := (arb.BodyA.deleted || arb.BodyB.deleted)
		disabled := !(arb.BodyA.Enabled || arb.BodyB.Enabled)
		if (ticks >= 1 && arb.state != arbiterStateCached) || deleted || disabled {
			arb.state = arbiterStateCached
			if arb.BodyA.CallbackHandler != nil {
				arb.BodyA.CallbackHandler.CollisionExit(arb)
			}
			if arb.BodyB.CallbackHandler != nil {
				arb.BodyB.CallbackHandler.CollisionExit(arb)
			}
		}
		if ticks > time.Duration(space.collisionPersistence) || deleted {
			delete(space.cachedArbiters, h)
			space.ArbiterBuffer = append(space.ArbiterBuffer, arb)
			c := arb.Contacts
			if c != nil {
				space.ContactBuffer = append(space.ContactBuffer, c)
			}
		}
	}

	slop := space.collisionSlop
	biasCoef := float32(1.0 - math.Pow(float64(space.collisionBias), float64(dt)))
	invdt := float32(1 / dt)
	for _, arb := range space.Arbiters {
		arb.preStep(invdt, slop, biasCoef)
	}

	for _, con := range space.Constraints {
		con.PreSolve()
		con.PreStep(dt)
	}

	ldamping := float32(math.Pow(float64(space.LinearDamping), float64(dt)))
	adamping := float32(math.Pow(float64(space.AngularDamping), float64(dt)))

	for _, body := range space.Bodies {
		if body.Enabled {
			if body.IgnoreGravity {
				body.UpdateVelocity(vect.Vector_Zero, ldamping, adamping, dt)
				continue
			}
			body.UpdateVelocity(space.Gravity, ldamping, adamping, dt)
		}
	}

	// dt_coef := float32(0)
	// if prev_dt != 0 {
	// 	dt_coef = dt / prev_dt
	// }

	for _, arb := range space.Arbiters {
		arb.applyCachedImpulse(dt)
	}

	for _, con := range space.Constraints {
		con.ApplyCachedImpulse(dt)
	}

	//fmt.Println("STEP")
	// start = time.Now()

	//fmt.Println("Arbiters", len(space.Arbiters), biasCoef, dt)
	//spew.Config.MaxDepth = 3
	//spew.Config.Indent = "\t"
	for i := 0; i < space.Iterations; i++ {
		for _, arb := range space.Arbiters {
			arb.applyImpulse()
			//spew.Dump(arb)
			//spew.Printf("%+v\n", arb)
		}

		for _, con := range space.Constraints {
			con.ApplyImpulse()
		}
	}

	//fmt.Println("####")
	//fmt.Println("")

	//MultiThreadGo()
	//for i:=0; i<8; i++ {
	//	<-done
	//}
	space.ApplyImpulsesTime = time.Since(start)

	for _, con := range space.Constraints {
		con.PostSolve()
	}

	for _, arb := range space.Arbiters {
		if arb.ShapeA.Body.CallbackHandler != nil {
			arb.ShapeA.Body.CallbackHandler.CollisionPostSolve(arb)
		}
		if arb.ShapeB.Body.CallbackHandler != nil {
			arb.ShapeB.Body.CallbackHandler.CollisionPostSolve(arb)
		}
	}

	if len(space.deleteBodies) > 0 {
		for _, body := range space.deleteBodies {
			space.removeBody(body)
		}
		space.deleteBodies = space.deleteBodies[0:0]
	}

	// stepEnd := time.Now()
	// space.StepTime = stepEnd.Sub(stepStart)
}

// var done = make(chan bool, 8)
// var start = make(chan bool, 8)

// func (space *Space) MultiThreadTest() {
// 	for {
// 		<-start
// 		for i := 0; i < space.Iterations/8; i++ {
// 			for _, arb := range space.Arbiters {
// 				if arb.ShapeA.IsSensor || arb.ShapeB.IsSensor {
// 					continue
// 				}
// 				arb.applyImpulse()
// 			}
// 		}
// 		done <- true
// 	}
// }

// func MultiThreadGo() {
// 	for i := 0; i < 8; i++ {
// 		start <- true
// 	}
// 	for i := 0; i < 8; i++ {
// 		<-done
// 	}
// }

// func PrintTree(node *Node) {
// 	if node != nil {
// 		fmt.Println("Parent:")
// 		fmt.Println(node.bb)
// 		fmt.Println("A:")
// 		PrintTree(node.A)
// 		fmt.Println("B:")
// 		PrintTree(node.B)
// 	}
// }

func (space *Space) Space() *Space {
	return space
}

func (space *Space) Query(obj Indexable, aabb AABB, fnc SpatialIndexQueryFunc) {
	space.activeShapes.Query(obj, aabb, fnc)
}

func (space *Space) QueryStatic(obj Indexable, aabb AABB, fnc SpatialIndexQueryFunc) {
	space.staticShapes.Query(obj, aabb, fnc)
}

func (space *Space) SpacePointQueryFirst(point vect.Vect, layers Layer, group int, checkSensors bool) (shape *Shape) {

	found := false
	pointFunc := func(a, b Indexable) {
		if found {
			return
		}
		shapeB := b.Shape()
		shapeA := a.Shape()
		if queryRejectShapes(shapeA, shapeB) {
			if !checkSensors && shapeB.IsSensor {
				return
			}
			contacts := space.pullContactBuffer()
			numContacts := collide(contacts, shapeA, shapeB)
			if numContacts <= 0 {
				space.pushContactBuffer(contacts)
				return
			}
			shape = shapeB
			found = true
		}
	}

	dot := NewCircle(vect.Vector_Zero, 0.5)
	dot.BB = dot.update(transform.NewTransform(point, 0))
	dot.Layer = layers
	dot.Group = group
	space.staticShapes.Query(dot, dot.AABB(), pointFunc)
	if found {
		return
	}
	space.activeShapes.Query(dot, dot.AABB(), pointFunc)

	return
}

func (space *Space) SpacePointQuery(point vect.Vect, layers Layer, group int, checkSensors bool) (shapes []*Shape) {

	pointFunc := func(a, b Indexable) {
		shapeB := b.Shape()
		shapeA := a.Shape()
		// log.Println(queryRejectShapes(shapeA, shapeB))
		if queryRejectShapes(shapeA, shapeB) {
			if !checkSensors && shapeB.IsSensor {
				return
			}
			contacts := space.pullContactBuffer()
			numContacts := collide(contacts, shapeA, shapeB)
			if numContacts <= 0 {
				space.pushContactBuffer(contacts)
				return
			}
			shapes = append(shapes, shapeB)
		}
	}

	dot := NewCircle(vect.Vector_Zero, 0.5)
	dot.BB = dot.update(transform.NewTransform(point, 0))
	dot.Layer = layers
	dot.Group = group
	space.staticShapes.Query(dot, dot.AABB(), pointFunc)
	space.activeShapes.Query(dot, dot.AABB(), pointFunc)

	return
}

/*
func (space *Space) SpacePointQuery(point vect.Vect, layers Layer, group Group, cpSpacePointQueryFunc func, void *data)
{
	struct PointQueryContext context = {point, layers, group, func, data};
	cpBB bb = cpBBNewForCircle(point, 0.0f);

	cpSpaceLock(space); {
    cpSpatialIndexQuery(space->activeShapes, &context, bb, (cpSpatialIndexQueryFunc)PointQuery, data);
    cpSpatialIndexQuery(space->staticShapes, &context, bb, (cpSpatialIndexQueryFunc)PointQuery, data);
	} cpSpaceUnlock(space, cpTrue);
}
*/
func (space *Space) ActiveBody(body *Body) error {
	if body.IsRogue() {
		return errors.New("Internal error: Attempting to activate a rouge body.")
	}

	space.Bodies = append(space.Bodies, body)

	for _, shape := range body.Shapes {
		space.staticShapes.Remove(shape)
		space.activeShapes.Insert(shape)
	}
	/*
		for _, arb := range body.Arbiters {
			bodyA := arb.BodyA
			if body == bodyA || bodyA.IsStatic() {

					int numContacts = arb->numContacts;
					cpContact *contacts = arb->contacts;

					// Restore contact values back to the space's contact buffer memory
					arb->contacts = cpContactBufferGetArray(space);
					memcpy(arb->contacts, contacts, numContacts*sizeof(cpContact));
					cpSpacePushContacts(space, numContacts);

					// Reinsert the arbiter into the arbiter cache
					arbHashID := hashPair(arb.BodyA.Hash()*20, arb.BodyB.Hash()*10)
					space.cachedArbiters[arbHashID] = arb

					// Update the arbiter's state
					arb.stamp = space.stamp
					space->arbiters = append(space->arbiters, arb)

					//cpfree(contacts);

			}
		}
	*/

	return nil
}

func (space *Space) ProcessComponents(dt float32) {

	sleep := math.IsInf(float64(space.sleepTimeThreshold), 0)
	// bodies := space.Bodies
	// _ = space.Bodies
	if sleep {
		dv := space.idleSpeedThreshold
		dvsq := float32(0)
		if dv == 0 {
			dvsq = dv * dv
		} else {
			dvsq = space.Gravity.LengthSqr() * dt * dt
		}

		for _, body := range space.Bodies {
			keThreshold := float32(0)
			if dvsq != 0 {
				keThreshold = body.m * dvsq
			}
			body.node.IdleTime = 0
			if body.KineticEnergy() <= keThreshold {
				body.node.IdleTime += dt
			}
		}
	}

	// for _, arb := range space.Arbiters {
	// 	a, b := arb.BodyA, arb.BodyB
	// 	_, _ = a, b
	// 	if sleep {
	// 		log.Println("sleep")
	// 	}
	// }
	/*
		// Awaken any sleeping bodies found and then push arbiters to the bodies' lists.
		cpArray *arbiters = space->arbiters;
		for(int i=0, count=arbiters->num; i<count; i++){
			cpArbiter *arb = (cpArbiter*)arbiters->arr[i];
			cpBody *a = arb->body_a, *b = arb->body_b;

			if(sleep){
				if((cpBodyIsRogue(b) && !cpBodyIsStatic(b)) || cpBodyIsSleeping(a)) cpBodyActivate(a);
				if((cpBodyIsRogue(a) && !cpBodyIsStatic(a)) || cpBodyIsSleeping(b)) cpBodyActivate(b);
			}

			cpBodyPushArbiter(a, arb);
			cpBodyPushArbiter(b, arb);
		}

		if(sleep){
			// Bodies should be held active if connected by a joint to a non-static rouge body.
			cpArray *constraints = space->constraints;
			for(int i=0; i<constraints->num; i++){
				cpConstraint *constraint = (cpConstraint *)constraints->arr[i];
				cpBody *a = constraint->a, *b = constraint->b;

				if(cpBodyIsRogue(b) && !cpBodyIsStatic(b)) cpBodyActivate(a);
				if(cpBodyIsRogue(a) && !cpBodyIsStatic(a)) cpBodyActivate(b);
			}

			// Generate components and deactivate sleeping ones
			for(int i=0; i<bodies->num;){
				cpBody *body = (cpBody*)bodies->arr[i];

				if(ComponentRoot(body) == NULL){
					// Body not in a component yet. Perform a DFS to flood fill mark
					// the component in the contact graph using this body as the root.
					FloodFillComponent(body, body);

					// Check if the component should be put to sleep.
					if(!ComponentActive(body, space->sleepTimeThreshold)){
						cpArrayPush(space->sleepingComponents, body);
						CP_BODY_FOREACH_COMPONENT(body, other) cpSpaceDeactivateBody(space, other);

						// cpSpaceDeactivateBody() removed the current body from the list.
						// Skip incrementing the index counter.
						continue;
					}
				}

				i++;

				// Only sleeping bodies retain their component node pointers.
				body->node.root = NULL;
				body->node.next = NULL;
			}
		}
	*/
}

// Creates an arbiter between the given shapes.
// If the shapes do not collide, arbiter.NumContact is zero.
func (space *Space) CreateArbiter(sa, sb *Shape) *Arbiter {

	var arb *Arbiter
	if len(space.ArbiterBuffer) > 0 {
		arb, space.ArbiterBuffer = space.ArbiterBuffer[len(space.ArbiterBuffer)-1], space.ArbiterBuffer[:len(space.ArbiterBuffer)-1]
	} else {
		for i := 0; i < ArbiterBufferSize/2; i++ {
			space.ArbiterBuffer = append(space.ArbiterBuffer, newArbiter())
		}
		arb = newArbiter()
	}
	//arb = newArbiter()

	if sa.ShapeType() > sb.ShapeType() {
		arb.ShapeA = sb
		arb.ShapeB = sa
	} else {
		arb.ShapeA = sa
		arb.ShapeB = sb
	}

	arb.BodyA = arb.ShapeA.Body
	arb.BodyB = arb.ShapeB.Body

	arb.Surface_vr = vect.Vect{}
	arb.stamp = 0
	//arb.nodeA = new(ArbiterEdge)
	//arb.nodeB = new(ArbiterEdge)
	arb.state = arbiterStateFirstColl
	arb.Contacts = nil
	arb.NumContacts = 0
	arb.e = 0
	arb.u = 0

	return arb
}

func spaceCollideShapes(a, b Indexable, null Data) {
	SpaceCollideShapes(a.Shape(), b.Shape(), a.Shape().space)
}

func SpaceCollideShapes(a, b *Shape, space *Space) {
	// log.Println(a.Group, b.Group, queryReject(a, b))
	// if a == nil || b == nil || a.Body == nil || b.Body == nil {
	// 	return
	// }

	if queryReject(a, b) {
		return
	}

	if a.ShapeType() > b.ShapeType() {
		a, b = b, a
	}

	//cpCollisionHandler *handler = cpSpaceLookupHandler(space, a->collision_type, b->collision_type);

	sensor := a.IsSensor || b.IsSensor
	//if(sensor && handler == &cpDefaultCollisionHandler) return;
	//if sensor {
	//	return
	//}

	// Narrow-phase collision detection.
	contacts := space.pullContactBuffer()

	numContacts := collide(contacts, a, b)
	// log.Println(numContacts)
	if numContacts <= 0 {
		space.pushContactBuffer(contacts)
		return // Shapes are not colliding.
	}

	contacts = contacts[:numContacts]

	// Get an arbiter from space->arbiterSet for the two shapes.
	// This is where the persistant contact magic comes from.

	arbHashID := newPair(a, b)

	var arb *Arbiter

	arb, exist := space.cachedArbiters[arbHashID]
	if !exist {
		arb = space.CreateArbiter(a, b)
	}

	var oldContacts []*Contact

	if arb.Contacts != nil {
		oldContacts = arb.Contacts
	}
	arb.update(a, b, contacts, numContacts)
	if oldContacts != nil {
		space.pushContactBuffer(oldContacts)
	}

	// if a == nil || b == nil || a.Body == nil || b.Body == nil {
	// 	return
	// }
	// log.Println(a.Body, b.Body.CallbackHandler, arb)

	space.cachedArbiters[arbHashID] = arb

	// Call the begin function first if it's the first step
	if arb.state == arbiterStateFirstColl {
		ignore := false
		if b.Body.CallBackCollision != nil {
			ignore = !b.Body.CallBackCollision(arb)
		}
		if a.Body.CallBackCollision != nil {
			ignore = ignore || !a.Body.CallBackCollision(arb)
		}

		if b.Body.CallbackHandler != nil {
			ignore = !b.Body.CallbackHandler.CollisionEnter(arb)
		}
		if a.Body.CallbackHandler != nil {
			ignore = ignore || !a.Body.CallbackHandler.CollisionEnter(arb)
		}
		if ignore {
			arb.Ignore() // permanently ignore the collision until separation
		}
	}

	preSolveResult := true

	// Ignore the arbiter if it has been flagged
	if arb.state != arbiterStateIgnore {
		// Call preSolve
		if arb.ShapeA.Body.CallbackHandler != nil {
			preSolveResult = arb.ShapeA.Body.CallbackHandler.CollisionPreSolve(arb)
		}
		if arb.ShapeB.Body.CallbackHandler != nil {
			preSolveResult = preSolveResult || arb.ShapeB.Body.CallbackHandler.CollisionPreSolve(arb)
		}
	} else {
		preSolveResult = false
	}

	if preSolveResult &&
		// Process, but don't add collisions for sensors.
		!sensor {
		space.Arbiters = append(space.Arbiters, arb)
	} else {
		//cpSpacePopContacts(space, numContacts);

		space.ContactBuffer = append(space.ContactBuffer, arb.Contacts)
		arb.Contacts = nil
		arb.NumContacts = 0

		// Normally arbiters are set as used after calling the post-solve callback.
		// However, post-solve callbacks are not called for sensors or arbiters rejected from pre-solve.
		if arb.state != arbiterStateIgnore {
			arb.state = arbiterStateNormal
		}
	}

	// Time stamp the arbiter so we know it was used recently.

	arb.stamp = space.stamp
}

func queryRejectShapes(a, b *Shape) bool {
	return a == b || (a.Group != 0 && a.Group == b.Group) || (a.Layer&b.Layer) == 0 || (a.Body != nil && !a.Body.Enabled) || (b.Body != nil && !b.Body.Enabled)
}

func queryReject(a, b *Shape) bool {
	// if a.Group == 0 && b.Group == 0 {
	// 	return true
	// }

	if a == b {
		return true
	}

	if a.Body == nil || b.Body == nil {
		return true
	}

	if a.Body == b.Body {
		return true
	}

	// if a.Group != b.Group {
	// 	return true
	// }

	if !a.Body.Enabled || !b.Body.Enabled {
		// log.Println("disabled")
		return true
	}

	return false
	//|| (a.Layer & b.Layer) != 0
	return a.Body == b.Body || (a.Group != 0 && a.Group == b.Group) || (a.Layer&b.Layer) == 0 || !a.Body.Enabled || !b.Body.Enabled || (math.IsInf(float64(a.Body.m), 0) && math.IsInf(float64(b.Body.m), 0)) || !TestOverlapPtr(&a.BB, &b.BB)
}

type RayCast struct {
	begin vect.Vect
	dir   vect.Vect
}

type RayCastHit struct {
	Distance float32
	Body     *Body
	Shape    *Shape
}

// const EPS = 0.00001

// const EPS = 0.1

func RayAgainstPolygon(c RayCast, poly *PolygonShape) bool {

	// poly.TestPoint(point)
	// log.Println("")
	// log.Println(c)

	// log.Println(poly.Axes, poly.TAxes)

	// log.Println(poly.ContainsVertPartial(c.begin, c.dir))

	for i, _ := range poly.TAxes {

		v1 := poly.TVerts[i]
		v2 := poly.TVerts[(i+1)%poly.NumVerts]

		intersect := vect.Intersection(c.begin, c.dir, v1, v2)
		if intersect {
			return true
		}
		continue

		// // log.Println(intersect, axis, v1, v2)

		// cosAngle := vect.Dot(c.dir, axis.N)

		// // if cosAngle < EPS && cosAngle >= -EPS {
		// // 	// log.Println("FALSE")
		// // 	return false
		// // }

		// t := -(vect.Dot(c.begin, axis.N) - axis.D) / cosAngle
		// log.Println(t)
		// // log.Println(t, axis, t > 1.1, t < -0.2)
		// // if t < -0.1 || t > 1 {
		// // 	log.Println("FALSE T", t)
		// // 	return false
		// // }
		// //check if point belongs to polygon line

		// point := vect.Add(c.begin, vect.Mult(c.dir, t))

		// polyDir := vect.Sub(v2, v1)

		// // log.Println(point, polyDir, v1, v2)

		// polyX := (point.X - v1.X) / polyDir.X
		// polyY := (point.Y - v1.Y) / polyDir.Y

		// if polyX >= -0.1 && polyX <= 1 {
		// 	log.Println("TRUE  X")
		// 	return true
		// }

		// if polyY >= -0.1 && polyY <= 1 {
		// 	log.Println("TRUE  Y")
		// 	return true
		// }
	}

	// log.Println("NOT FOUND")
	return false
}

func pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func sqr(x float32) float32 {
	return x * x
}

// // SegmentCircleIntersection return points of intersection between a circle and
// // a line segment. The Boolean intersects returns true if one or
// // more solutions exist. If only one solution exists,
// // x1 == x2 and y1 == y2.
// // s1x and s1y are coordinates for one end point of the segment, and
// // s2x and s2y are coordinates for the other end of the segment.
// // cx and cy are the coordinates of the center of the circle and
// // r is the radius of the circle.
// func SegmentCircleIntersection(s1x, s1y, s2x, s2y, cx, cy, r float32) (x1, y1, x2, y2 float32, intersects bool) {
// 	log.Println(s1x, s1y, ":", s2x, s2y)
// 	log.Println(cx, cy, r)

// 	// (n-et) and (m-dt) are expressions for the x and y coordinates
// 	// of a parameterized line in coordinates whose origin is the
// 	// center of the circle.
// 	// When t = 0, (n-et) == s1x - cx and (m-dt) == s1y - cy
// 	// When t = 1, (n-et) == s2x - cx and (m-dt) == s2y - cy.
// 	n := s2x - cx
// 	m := s2y - cy

// 	e := s2x - s1x
// 	d := s2y - s1y

// 	// lineFunc checks if the  t parameter is in the segment and if so
// 	// calculates the line point in the unshifted coordinates (adds back
// 	// cx and cy.
// 	lineFunc := func(t float32) (x, y float32, inBounds bool) {
// 		inBounds = t >= 0 && t <= 1 // Check bounds on closed segment
// 		// To check bounds for an open segment use t > 0 && t < 1
// 		if inBounds { // Calc coords for point in segment
// 			x = n - e*t + cx
// 			y = m - d*t + cy
// 		}
// 		return
// 	}

// 	// Since we want the points on the line distance r from the origin,
// 	// (n-et)(n-et) + (m-dt)(m-dt) = rr.
// 	// Expanding and collecting terms yeilds the following quadratic equation:
// 	A, B, C := e*e+d*d, -2*(e*n+m*d), n*n+m*m-r*r

// 	D := B*B - 4*A*C // discriminant of quadratic
// 	if D < 0 {
// 		return // No solution
// 	}
// 	D = float32(math.Sqrt(float64(D)))

// 	var p1In, p2In bool
// 	x1, y1, p1In = lineFunc((-B + D) / (2 * A)) // First root
// 	if D == 0 {
// 		intersects = p1In
// 		x2, y2 = x1, y1
// 		return // Only possible solution, quadratic has one root.
// 	}

// 	x2, y2, p2In = lineFunc((-B - D) / (2 * A)) // Second root

// 	intersects = p1In || p2In
// 	if p1In == false { // Only x2, y2 may be valid solutions
// 		x1, y1 = x2, y2
// 	} else if p2In == false { // Only x1, y1 are valid solutions
// 		x2, y2 = x1, y1
// 	}
// 	return
// }

//RayAgainstCircle... mathematic black magic!
func RayAgainstCircle(ray RayCast, circle *CircleShape) bool {

	// x1, x2, y1, y2, ok := SegmentCircleIntersection(ray.begin.X, ray.begin.Y, ray.dir.X, ray.dir.Y, circle.Tc.X, circle.Tc.Y, circle.Radius)
	// log.Println(x1, x2, y1, y2, ok)

	// return ok

	/////////////////////////////

	// x1 := ray.begin.X
	// y1 := ray.begin.Y
	// x2 := ray.dir.X
	// y2 := ray.dir.Y

	// xC := circle.Tc.X
	// yC := circle.Tc.Y
	// R := circle.Radius

	// x1 -= xC
	// y1 -= yC
	// x2 -= xC
	// y2 -= yC

	// dx := x2 - x1
	// dy := y2 - y1

	// a := dx*dx + dy*dy
	// b := 2 * (x1*dx + y1*dy)
	// c := x1*x1 + y1*y1 - R*R

	// if -b < 0 {
	// 	return c < 0
	// }
	// if -b < 2*a {
	// 	return 4*a*c-b*b < 0
	// }
	// return a+b+c < 0

	////////////////////////////

	x0 := circle.Tc.X
	y0 := circle.Tc.Y
	r := circle.Radius

	x1 := ray.begin.X
	y1 := ray.begin.Y
	x2 := ray.dir.X
	y2 := ray.dir.Y

	var dx01 = x1 - x0
	var dy01 = y1 - y0
	var dx12 = x2 - x1
	var dy12 = y2 - y1

	var a = sqr(dx12) + sqr(dy12)
	var k = dx01*dx12 + dy01*dy12
	var c = sqr(dx01) + sqr(dy01) - sqr(r)

	var d1 = sqr(k) - a*c
	if d1 >= 0 && k < 0 {
		return true
	}

	return false

	////////////////////////////

	// // fromRayToCircle := vect.Sub(cast.begin, circle.Tc)
	// // a := cast.dir.LengthSqr()
	// // b := 2.0 * vect.Dot(fromRayToCircle, cast.dir)
	// // c := vect.Dot(fromRayToCircle, fromRayToCircle) - circle.Radius*circle.Radius

	// // D := b*b - 4.0*a*c

	// // if D < 0.0 {
	// // 	return false
	// // }
	// // D = float32(math.Sqrt(float64(D)))
	// // t1 := (-b - D) / (2.0 * a)
	// // t2 := (-b + D) / (2.0 * a)

	// // if (t1 >= 0.0 && t1 <= 1.0) || (t2 >= 0.0 && t2 <= 1.0) {
	// // 	return true
	// // }
	// // return false
}

func Distance(x0, y0, x1, y1 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x0)-float64(x1), 2) + math.Pow(float64(y0)-float64(y1), 2)))
}

func (space *Space) RayCastAll(begin vect.Vect, direction vect.Vect, group int, ignoreBody *Body) (hits []*RayCastHit) {

	rayCast := RayCast{
		begin: begin,
		dir:   direction,
	}

	length := Distance(begin.X, begin.Y, direction.X, direction.Y)
	// maxLength := length * 1 //crutch

	// space.staticShapes.Contains(obj)

	// dot := NewBox(begin, 0.1, length)
	// dot.Group = group
	// space.
	// dot.BB = dot.update(transform.NewTransform(point, 0))
	// dot.Layer = layers
	// dot.Group = group
	// space.staticShapes.Query(dot, dot.AABB(), pointFunc)
	// space.activeShapes.Query(dot, dot.AABB(), pointFunc)

	for _, body := range space.Bodies {
		if body == ignoreBody {
			continue
		}

		for _, shape := range body.Shapes {

			if shape.Group != group {
				continue
			}

			pos := body.Position()
			dist := Distance(pos.X, pos.Y, begin.X, begin.Y)

			if dist > length || dist <= 0 {
				// log.Println("object too far", body.UserData)
				continue
			}

			var hit bool

			switch shape.ShapeType() {
			case ShapeType_Polygon:
				hit = RayAgainstPolygon(rayCast, shape.GetAsPolygon())
			case ShapeType_Circle:
				hit = RayAgainstCircle(rayCast, shape.GetAsCircle())
			case ShapeType_Box:
				hit = RayAgainstPolygon(rayCast, shape.GetAsBox().Polygon)
			default:
				log.Printf("WARNING: shape type `%s` not work", shape.ShapeType())
			}

			if hit {
				// log.Println("HIT:", dist, body.UserData)
				hits = append(hits, &RayCastHit{
					Distance: dist,
					Body:     body,
					Shape:    shape,
				})
				continue
			}
		}
	}

	// log.Println("HITS:", len(hits))

	return
}

func (space *Space) AddBody(body *Body) *Body {
	if body.space != nil {
		println("This body is already added to a space and cannot be added to another.")
		return body
	}

	body.space = space
	// if !body.IsStatic() {
	space.Bodies = append(space.Bodies, body)
	// }

	for _, shape := range body.Shapes {
		if shape.space == nil {
			space.AddShape(shape)
		}
	}

	return body
}

func (space *Space) AddShape(shape *Shape) *Shape {
	if shape.space != nil {
		println("This shape is already added to a space and cannot be added to another.")
		return shape
	}

	shape.space = space
	shape.Update()
	if shape.Body.IsStatic() {
		space.staticShapes.Insert(shape)
	} else {
		space.activeShapes.Insert(shape)
	}

	return shape
}

func (space *Space) AddConstraint(constraint Constraint) Constraint {
	con := constraint.Constraint()
	if con.space != nil {
		panic("This shape is already added to a space and cannot be added to another.")
	}

	con.BodyA.BodyActivate()
	con.BodyB.BodyActivate()
	space.Constraints = append(space.Constraints, constraint)

	// Push onto the heads of the bodies' constraint lists
	//cpBody *a = constraint->a, *b = constraint->b;
	//constraint->next_a = a->constraintList; a->constraintList = constraint;
	//constraint->next_b = b->constraintList; b->constraintList = constraint;
	con.space = space

	return constraint
}

func (space *Space) RemoveConstraint(constraint Constraint) {
	con := constraint.Constraint()
	if con.space == nil {
		panic("Cannot remove a constraint that was not added to the space. (Removed twice maybe?)")
	}

	con.BodyA.BodyActivate()
	con.BodyB.BodyActivate()

	for i, c := range space.Constraints {
		if constraint == c {
			space.Constraints[i], space.Constraints = space.Constraints[len(space.Constraints)-1], space.Constraints[:len(space.Constraints)-1]
			break
		}
	}

	//cpBodyRemoveConstraint(constraint->a, constraint);
	//cpBodyRemoveConstraint(constraint->b, constraint);
	con.space = nil
	con.BodyA = nil
	con.BodyB = nil
}

func (space *Space) removeBody(body *Body) {
	for _, shape := range body.Shapes {
		space.RemoveShape(shape)
	}
	body.space = nil
	body.Shapes = nil
	body.UserData = nil
	body.CallbackHandler = nil
	body.UpdateVelocityFunc = nil
	body.UpdatePositionFunc = nil
}

func (space *Space) RemoveBody(body *Body) {
	if body == nil {
		return
	}
	body.BodyActivate()
	for i, pbody := range space.Bodies {
		if pbody == body {
			space.Bodies[i], space.Bodies = space.Bodies[len(space.Bodies)-1], space.Bodies[:len(space.Bodies)-1]
			break
		}
	}
	body.deleted = true
	space.deleteBodies = append(space.deleteBodies, body)
}

func (space *Space) RemoveShape(shape *Shape) {
	shape.space = nil
	if shape.Body.IsStatic() {
		space.staticShapes.Remove(shape)
	} else {
		space.activeShapes.Remove(shape)
	}
	shape.Body = nil
	shape.UserData = nil
	shape.ShapeClass = nil
}

func (space *Space) pullContactBuffer() (contacts []*Contact) {
	if len(space.ContactBuffer) > 0 {
		contacts, space.ContactBuffer = space.ContactBuffer[len(space.ContactBuffer)-1], space.ContactBuffer[:len(space.ContactBuffer)-1]
	} else {
		for i := 0; i < ContactBufferSize/2; i++ {
			ccs := make([]*Contact, MaxPoints)

			for i := 0; i < MaxPoints; i++ {
				ccs[i] = &Contact{}
			}
			space.ContactBuffer = append(space.ContactBuffer, ccs)
		}
		contacts, space.ContactBuffer = space.ContactBuffer[len(space.ContactBuffer)-1], space.ContactBuffer[:len(space.ContactBuffer)-1]
	}
	return
}

func (space *Space) pushContactBuffer(contacts []*Contact) {
	space.ContactBuffer = append(space.ContactBuffer, contacts)
}
