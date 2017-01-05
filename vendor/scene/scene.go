package scene

import (
	"engine"
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v2"
)

var Dir = "assets/scenes"

//Scene structure
type Scene struct {
	Name    string
	Objects []*engine.Object
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
		log.Fatalln("failed load scene:", err)
	}

	for _, o := range s.Objects {
		o.P.Static = true
		log.Println(o.P.Pos)
		o.Create()
	}
}
