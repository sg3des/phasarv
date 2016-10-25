package scene

import (
	"engine"
	"io/ioutil"
	"log"
	"param"
	"path"

	"gopkg.in/yaml.v2"
)

var Dir = "assets/scenes"

//Scene structure
type Scene struct {
	Name    string
	Objects []struct {
		Object param.Object
		HP     float32 `yaml:"hp"`
	}

	Shaders []string
}

//read yaml file and parse to Scene structure
func read(name string) (Scene, error) {
	var s Scene

	data, err := ioutil.ReadFile(path.Join(Dir, name+".yaml"))
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(data, &s)

	return s, err
}

//Load scene from file
func Load(name string) {
	s, err := read(name)
	if err != nil {
		log.Fatalln("scene not found", err)
	}

	for _, o := range s.Objects {
		engine.NewStaticObject(o.Object)
	}
}
