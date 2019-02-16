package scenes

import (
	"engine"
	"io/ioutil"
	"log"
	"path"
	"phys"
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"

	"gopkg.in/yaml.v2"
)

var (
	Dir    = "assets/scenes"
	Camera *fizzle.YawPitchCamera
)

//Scene structure
type Scene struct {
	Name    string
	Objects []*engine.Object
}

//read yaml file and parse to Scene structure
func read(name string) (*Scene, error) {
	s := &Scene{Name: name}

	data, err := ioutil.ReadFile(path.Join(Dir, name+".yaml"))
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(data, s)

	return s, err
}

//Load scene from file
func Load(name string) *Scene {
	s, err := read(name)
	if err != nil {
		log.Fatalln("failed load scene:", err)
	}

	for _, o := range s.Objects {
		// o.Name = fmt.Sprintf("%s-%d", o.Name, i)
		o.P.Static = true
		if o.PI != nil {
			o.PI.Group = phys.GROUP_STATIC
		}
		o.Create()
	}

	return s
}

func (s *Scene) Close() {
	if s == nil {
		return
	}

	for _, o := range s.Objects {
		o.Remove()
	}
}

func InitEnvironment() {
	sun := &render.Light{
		Direct:     true,
		Pos:        mgl32.Vec3{-3, 3, 10},
		Dir:        mgl32.Vec3{5, -5, -10},
		Strength:   0.8, //0.8,
		Specular:   0.2,
		ShadowSize: 4096,
		Distance:   100,
	}
	sun.Create()

	backlight := &render.Light{
		Direct:     true,
		Pos:        mgl32.Vec3{0, 0, 5},
		Dir:        mgl32.Vec3{0, 0, -5},
		Strength:   0.2,
		Specular:   0.1,
		ShadowSize: 1,
		Distance:   15,
	}
	backlight.Create()

	Camera = render.NewCamera(mgl32.Vec3{0, 0, 40})
	Camera.LookAtDirect(mgl32.Vec3{0, 0, 0})
}
